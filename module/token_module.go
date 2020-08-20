package module

import (
	"fmt"
	"github.com/adschain/log"
	"github.com/wlcy/tradehome-service/common/mysql"
	"github.com/wlcy/tradehome-service/entity"
	"github.com/wlcy/tradehome-service/util"
)

func QueryTokenInfo(address string) (*entity.TokenInfo, error) {
	strSQL := fmt.Sprintf(`
		select id, address, full_name, short_name, circulation, token_precision, description, website_url, logo_url
		from market_token_info where address='%v'`, address)
	log.Info(strSQL)
	dataPtr, err := mysql.QueryTableData(strSQL)
	if err != nil {
		log.Errorf(err, "QueryTokenInfo error sql:%v", strSQL)
		return nil, err
	}

	item := &entity.TokenInfo{}
	for dataPtr.NextT() {
		item.Id = util.ConvertDBValueToInt64(dataPtr.GetField("id"))
		item.Address = dataPtr.GetField("address")
		item.FullName = dataPtr.GetField("full_name")
		item.ShortName = dataPtr.GetField("short_name")
		item.Circulation = util.ConvertDBValueToInt64(dataPtr.GetField("circulation"))
		precision := util.ConvertDBValueToInt(dataPtr.GetField("token_precision"))
		item.Precision = precision
		item.Description = dataPtr.GetField("description")
		item.WebsiteUrl = dataPtr.GetField("website_url")
		item.LogoUrl = dataPtr.GetField("logo_url")
	}

	return item, nil
}
