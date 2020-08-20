package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/util"
	"time"
)

// ----------------------------stm_divide-------------------------------
//查询数据库中的计算结果
func queryPreparedAndSendReward() ([]*UserRewardInfo, error) {
	now := fmt.Sprintf("%v 00:00:00", time.Now().Format(util.DATEFORMAT))
	strSQL := fmt.Sprintf(`
	SELECT id, address,pledged, total_pledged,total_bonus,reward,trx_hash,trx_result 
    FROM stm_divide 
    where trx_result='' and trx_hash='' and c_time>='%v'`, now)
	log.Infof(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Infof("queryTokenExtInfo error :[%v]\n", err)
		return nil, err
	}
	if dataPtr == nil {
		log.Infof("dataPtr is nil")
		return nil, errors.New("dataPtr is nil")
	}

	var userRewards = make([]*UserRewardInfo, 0)

	for dataPtr.NextT() {
		reward := &UserRewardInfo{}
		reward.ID = util.ConvertDBValueToInt64(dataPtr.GetField("id"))
		reward.UserAddr = dataPtr.GetField("address")
		reward.Pledged = util.ConvertDBValueToInt64(dataPtr.GetField("pledged"))
		reward.TotalPledged = util.ConvertDBValueToInt64(dataPtr.GetField("total_pledged"))
		reward.TotalBonus = util.ConvertDBValueToInt64(dataPtr.GetField("total_bonus"))
		reward.Bonus = util.ConvertDBValueToInt64(dataPtr.GetField("reward"))
		reward.TrxHash = dataPtr.GetField("trx_hash")
		reward.TrxResult = dataPtr.GetField("trx_result")
		userRewards = append(userRewards, reward)
	}
	return userRewards, nil
}

//仅插入奖励计算结果======================
func InsertRewardStandby(rewardInfos []*UserRewardInfo) error {
	if len(rewardInfos) == 0 {
		return errors.New("insert db with no data")
	}
	txn, err := mysql.OpenDataBaseWriteTransaction()
	if err != nil {
		log.Infof("open database transaction err:%v", err)
		return errors.New("open database transaction err")
	}
	if nil == txn {
		log.Infof("open database instance is nil")
		return fmt.Errorf("open database instance is nil")
	}

	defer func() {
		if recErr := recover(); recErr != nil || err != nil {
			log.Infof("recErr:[%v], err:[%v]", recErr, err)
			txn.Rollback()
		} else {
			txn.Commit()
		}
	}()

	err = insertBonusBatch(rewardInfos, txn)
	if err != nil {
		log.Infof("db exchange price  insert err: %#v", err)
		return err
	}
	return nil
}

func insertBonusBatch(resultInfos []*UserRewardInfo, tnx *sql.Tx) (err error) {
	ts := time.Now()

	strSQL := fmt.Sprintf(`INSERT INTO stm_divide (address,pledged, total_pledged,total_bonus,reward,trx_hash,trx_result)
	 values(?,?,?,?,?,?,?)`)
	log.Infof("update--stm_divide-started-------[%v]-----", strSQL)

	stmt, err := tnx.Prepare(strSQL)
	if err != nil || stmt == nil {
		fmt.Println(err)
		return errors.New("tnx prepare err")
	}
	if nil != stmt {
		defer stmt.Close()
	}

	insertCnt := 0
	errCnt := 0
	for _, data := range resultInfos {
		log.Infof("insert data:[%#v]", data)
		_, err := stmt.Exec(data.UserAddr, data.Pledged, data.TotalPledged, data.TotalBonus, data.Bonus, data.TrxHash, data.TrxResult)
		if err != nil {
			log.Infof("prepare data err:[%v]", err)
			errCnt++
		} else {
			insertCnt++
		}

	}
	log.Infof("store stm_divide record OK, cost:%v, insertCnt:%v, errCnt:%v, total source:%v\n", time.Since(ts), insertCnt, errCnt, len(resultInfos))
	return nil
}

// 批量更新数据库奖励发放结果================
func UpdateRewardResult(rewardInfos []*UserRewardInfo) error {
	if len(rewardInfos) == 0 {
		return errors.New("insert db with no data")
	}
	txn, err := mysql.OpenDataBaseWriteTransaction()
	if err != nil {
		log.Infof("open database transaction err:%v", err)
		return errors.New("open database transaction err")
	}
	if nil == txn {
		log.Infof("open database instance is nil")
		return fmt.Errorf("open database instance is nil")
	}

	defer func() {
		if recErr := recover(); recErr != nil || err != nil {
			log.Infof("recErr:[%v], err:[%v]", recErr, err)
			txn.Rollback()
		} else {
			txn.Commit()
		}
	}()

	err = updateBonusBatch(rewardInfos, txn)
	if err != nil {
		log.Infof("db update bonus err: %#v", err)
		return err
	}
	return nil
}

func updateBonusBatch(resultInfos []*UserRewardInfo, tnx *sql.Tx) (err error) {
	ts := time.Now()
	strSQL := fmt.Sprintf(`update stm_divide set trx_hash=?,trx_result=? 
    where id=?`)
	log.Infof("update--stm_divide-started-------[%v]-----", strSQL)

	stmt, err := tnx.Prepare(strSQL)
	if err != nil || stmt == nil {
		fmt.Println(err)
		return errors.New("tnx prepare err")
	}
	if nil != stmt {
		defer stmt.Close()
	}

	insertCnt := 0
	errCnt := 0
	for _, data := range resultInfos {
		log.Infof("update data:[%#v]", data)
		_, err := stmt.Exec(data.TrxHash, data.TrxResult, data.ID)
		if err != nil {
			log.Infof("prepare data err:[%v]", err)
			errCnt++
		} else {
			insertCnt++
		}

	}
	log.Infof("store stm_divide record OK, cost:%v, updateCnt:%v, errCnt:%v, total source:%v\n", time.Since(ts), insertCnt, errCnt, len(resultInfos))
	return nil
}

//----------------------------- stm_pool -------------------------------------------------

//查询数据库中的计算结果
func queryRewardPoolInfo() ([]*DividePoolInfo, error) {
	now := fmt.Sprintf("%v 00:00:00", time.Now().Format(util.DATEFORMAT))
	end := fmt.Sprintf("%v 23:59:59", time.Now().Format(util.DATEFORMAT))
	strSQL := fmt.Sprintf(`
	SELECT caller_address,pool_address, total_amount,trx_hash,trx_result 
    FROM stm_pool
    where c_time>='%v' and c_time<'%v'`, now, end)
	log.Infof(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Infof("queryRewardPoolInfo error :[%v]\n", err)
		return nil, err
	}
	if dataPtr == nil {
		log.Infof("dataPtr is nil")
		return nil, errors.New("dataPtr is nil")
	}

	var poolInfos = make([]*DividePoolInfo, 0)

	for dataPtr.NextT() {
		poolInfo := &DividePoolInfo{}
		poolInfo.ID = util.ConvertDBValueToInt64(dataPtr.GetField("id"))
		poolInfo.CallerAddress = dataPtr.GetField("caller_address")
		poolInfo.PoolAddress = dataPtr.GetField("pool_address")
		poolInfo.TotalAmount = util.ConvertDBValueToInt64(dataPtr.GetField("total_amount"))
		poolInfo.TrxHash = dataPtr.GetField("trx_hash")
		poolInfo.TrxResult = dataPtr.GetField("trx_result")
		poolInfos = append(poolInfos, poolInfo)
	}
	return poolInfos, nil
}

//InsertPoolReward divide 调用结果，记录奖池金额======================
func InsertPoolReward(poolInfo *DividePoolInfo) error {
	txn, err := mysql.OpenDataBaseWriteTransaction()
	if err != nil {
		log.Infof("open database transaction err:%v", err)
		return errors.New("open database transaction err")
	}
	if nil == txn {
		log.Infof("open database instance is nil")
		return fmt.Errorf("open database instance is nil")
	}

	defer func() {
		if recErr := recover(); recErr != nil || err != nil {
			log.Infof("recErr:[%v], err:[%v]", recErr, err)
			txn.Rollback()
		} else {
			txn.Commit()
		}
	}()

	strSQL := fmt.Sprintf(`INSERT INTO stm_pool (caller_address,pool_address, total_amount,trx_hash,trx_result)
	 values(?,?,?,?,?)`)
	log.Infof("update--stm_pool-started-------[%v]-----", strSQL)

	stmt, err := txn.Prepare(strSQL)
	if err != nil || stmt == nil {
		fmt.Println(err)
		return errors.New("tnx prepare err")
	}
	if nil != stmt {
		defer stmt.Close()
	}

	_, err = stmt.Exec(poolInfo.CallerAddress, poolInfo.PoolAddress, poolInfo.TotalAmount, poolInfo.TrxHash, poolInfo.TrxResult)
	if err != nil {
		log.Infof("db stm_pool insert err: %#v", err)
		return err
	}
	log.Infof("store stm_pool record OK \n")

	return nil
}
