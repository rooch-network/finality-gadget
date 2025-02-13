package db

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/babylonlabs-io/finality-gadget/types"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

type BBoltHandler struct {
	db     *bolt.DB
	logger *zap.Logger
}

var _ IDatabaseHandler = &BBoltHandler{}

const (
	blocksBucket          = "blocks"
	blockHeightsBucket    = "block_heights"
	indexerBucket         = "indexer"
	earliestBlockKey      = "earliest"
	latestBlockKey        = "latest"
	activatedTimestampKey = "activated_timestamp"
)

//////////////////////////////
// CONSTRUCTOR
//////////////////////////////

func NewBBoltHandler(path string, logger *zap.Logger) (*BBoltHandler, error) {
	// 0600 = read/write permission for owner only
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Error("Error opening DB", zap.Error(err))
		return nil, err
	}

	return &BBoltHandler{
		db:     db,
		logger: logger,
	}, nil
}

//////////////////////////////
// METHODS
//////////////////////////////

func (bb *BBoltHandler) CreateInitialSchema() error {
	bb.logger.Info("Initialising DB...")
	return bb.db.Update(func(tx *bolt.Tx) error {
		buckets := []string{blocksBucket, blockHeightsBucket, indexerBucket}
		for _, bucket := range buckets {
			if err := bb.tryCreateBucket(tx, bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (bb *BBoltHandler) InsertBlocks(blocks []*types.Block) error {
	if len(blocks) == 0 {
		return nil
	}

	bb.logger.Info("Batch inserting blocks to DB", zap.Int("count", len(blocks)))

	// Single transaction for all operations
	return bb.db.Update(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte(blocksBucket))
		heightsBucket := tx.Bucket([]byte(blockHeightsBucket))
		indexBucket := tx.Bucket([]byte(indexerBucket))

		var minHeight, maxHeight uint64 = math.MaxUint64, 0

		// Insert all blocks
		for _, block := range blocks {
			// Update min/max heights
			if block.BlockHeight < minHeight {
				bb.logger.Debug("Setting min height to", zap.Uint64("min_height", block.BlockHeight))
				minHeight = block.BlockHeight
			}
			if block.BlockHeight > maxHeight {
				bb.logger.Debug("Setting max height to", zap.Uint64("max_height", block.BlockHeight))
				maxHeight = block.BlockHeight
			}

			// Store block data
			blockBytes, err := json.Marshal(block)
			if err != nil {
				bb.logger.Error("Error inserting block", zap.Error(err))
				return err
			}
			bb.logger.Debug("Inserting block to db", zap.Uint64("block_height", block.BlockHeight), zap.String("block_hash", block.BlockHash))
			if err := blocksBucket.Put(bb.itob(block.BlockHeight), blockBytes); err != nil {
				bb.logger.Error("Error inserting block to db", zap.Error(err))
				return err
			}

			// Store height mapping
			bb.logger.Debug("Inserting height mapping to db", zap.String("block_hash", block.BlockHash), zap.Uint64("block_height", block.BlockHeight))
			if err := heightsBucket.Put([]byte(block.BlockHash), bb.itob(block.BlockHeight)); err != nil {
				bb.logger.Error("Error inserting height mapping", zap.Error(err))
				return err
			}
		}

		// Update earliest block if needed
		earliestBytes := indexBucket.Get([]byte(earliestBlockKey))
		if earliestBytes == nil {
			bb.logger.Debug("Updating earliest block in db", zap.Uint64("block_height", minHeight))
			if err := indexBucket.Put([]byte(earliestBlockKey), bb.itob(minHeight)); err != nil {
				bb.logger.Error("Error inserting earliest block", zap.Error(err))
				return err
			}
		}

		// Update latest block if needed
		latestBytes := indexBucket.Get([]byte(latestBlockKey))
		var currentLatest uint64
		if latestBytes != nil {
			currentLatest = bb.btoi(latestBytes)
		}
		if maxHeight > currentLatest {
			bb.logger.Debug("Updating latest block in db", zap.Uint64("block_height", maxHeight))
			if err := indexBucket.Put([]byte(latestBlockKey), bb.itob(maxHeight)); err != nil {
				bb.logger.Error("Error inserting latest block", zap.Error(err))
				return err
			}
		}

		return nil
	})
}

func (bb *BBoltHandler) GetBlockByHeight(height uint64) (*types.Block, error) {
	var block types.Block
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		v := b.Get(bb.itob(height))
		if v == nil {
			return types.ErrBlockNotFound
		}
		return json.Unmarshal(v, &block)
	})
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (bb *BBoltHandler) GetBlockByHash(hash string) (*types.Block, error) {
	// Fetch block number corresponding to hash
	var blockHeight uint64
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockHeightsBucket))
		v := b.Get([]byte(hash))
		if v == nil {
			return types.ErrBlockNotFound
		}
		blockHeight = bb.btoi(v)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bb.GetBlockByHeight(blockHeight)
}

func (bb *BBoltHandler) QueryIsBlockFinalizedByHeight(height uint64) (bool, error) {
	_, err := bb.GetBlockByHeight(height)
	if err != nil {
		if errors.Is(err, types.ErrBlockNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (bb *BBoltHandler) QueryIsBlockFinalizedByHash(hash string) (bool, error) {
	// Fetch block number by hash
	var blockHeight uint64
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockHeightsBucket))
		res := b.Get([]byte(hash))
		if len(res) == 0 {
			return types.ErrBlockNotFound
		}
		blockHeight = bb.btoi(res)
		return nil
	})
	if err != nil {
		if errors.Is(err, types.ErrBlockNotFound) {
			return false, nil
		}
		return false, err
	}

	return bb.QueryIsBlockFinalizedByHeight(blockHeight)
}

func (bb *BBoltHandler) QueryEarliestFinalizedBlock() (*types.Block, error) {
	var earliestBlockHeight uint64
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(indexerBucket))
		v := b.Get([]byte(earliestBlockKey))
		if v == nil {
			return types.ErrBlockNotFound
		}
		earliestBlockHeight = bb.btoi(v)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return bb.GetBlockByHeight(earliestBlockHeight)
}

func (bb *BBoltHandler) QueryLatestFinalizedBlock() (*types.Block, error) {
	var latestBlockHeight uint64

	// Fetch latest block height
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(indexerBucket))
		v := b.Get([]byte(latestBlockKey))
		if v == nil {
			return types.ErrBlockNotFound
		}
		latestBlockHeight = bb.btoi(v)
		return nil
	})
	if err != nil {
		// If no latest block has been stored yet, return nil
		if errors.Is(err, types.ErrBlockNotFound) {
			return nil, nil
		}
		bb.logger.Error("Error getting latest block", zap.Error(err))
		return nil, err
	}

	// Fetch latest block by height
	return bb.GetBlockByHeight(latestBlockHeight)
}

func (bb *BBoltHandler) GetActivatedTimestamp() (uint64, error) {
	var timestamp uint64
	err := bb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(indexerBucket))
		v := b.Get([]byte(activatedTimestampKey))
		if v == nil {
			return types.ErrActivatedTimestampNotFound
		}
		timestamp = bb.btoi(v)
		return nil
	})
	if err != nil {
		return math.MaxUint64, err
	}
	return timestamp, nil
}

func (bb *BBoltHandler) SaveActivatedTimestamp(timestamp uint64) error {
	return bb.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(indexerBucket))
		return b.Put([]byte(activatedTimestampKey), bb.itob(timestamp))
	})
}

func (bb *BBoltHandler) Close() error {
	bb.logger.Info("Closing DB...")
	return bb.db.Close()
}

//////////////////////////////
// INTERNAL
//////////////////////////////

func (bb *BBoltHandler) tryCreateBucket(tx *bolt.Tx, bucketName string) error {
	_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
	if err != nil {
		bb.logger.Error("Error creating bucket", zap.Error(err))
	}
	return err
}

func (bb *BBoltHandler) itob(v uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, v)
	if err != nil {
		bb.logger.Fatal("Error writing to buffer", zap.Error(err))
	}
	return buf.Bytes()
}

func (bb *BBoltHandler) btoi(b []byte) uint64 {
	var v uint64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &v)
	if err != nil {
		bb.logger.Fatal("Error reading from buffer", zap.Error(err))
	}
	return v
}
