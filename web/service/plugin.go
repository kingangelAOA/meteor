package service

import (
	"common"
	"common/bindmodels"
	"encoding/json"
	"fmt"
	"plugin/shared"
	"web/configs"
	"web/db"
	"web/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InitPlugin() {
	ps, err := GetPluginList()
	if err != nil {
		panic(err)
	}
	for _, p := range ps {
		pf := shared.NewPluginFactory(p.FileID, p.Code, p.Language, *configs.PluginRootPath)
		pc, err := pf.GetPluginClientFromDB(db.DownloadPluginFile)
		if err != nil {
			fmt.Println(err.Error())
		}
		if pc == nil {
			continue
		}
		if err := shared.PM.UpdatePlugin(p.ID, pc); err != nil {
			panic(err)
		}
	}
}

func CloseAllPlugin() error {
	return shared.PM.ClearAllPlugin()
}

func GetPlugins(q bindmodels.Query) ([]models.PluginVO, error) {
	var pvs []models.PluginVO
	ps, err := db.GetPlugins(q)
	if err != nil {
		return pvs, nil
	}
	for i := 0; i < len(ps); i++ {
		pv := convertToPluginVO(ps[i])
		pv.Status = CheckPluginHealth(pv.ID)
		pv.Index = i + 1 + int(q.Offset)
		pv.VersionNum = shared.PM.GetPluginVersionNum(pv.ID)
		pvs = append(pvs, pv)
	}
	return pvs, nil
}

func GetPluginList() ([]models.PluginVO, error) {
	var pvs []models.PluginVO
	ps, err := db.GetPluginList()
	if err != nil {
		return pvs, nil
	}
	for _, p := range ps {
		pvs = append(pvs, models.PluginVO{
			ID:       p.ID.Hex(),
			FileID:   p.FileID.Hex(),
			Name:     p.Name,
			Language: p.Language,
		})
	}
	return pvs, nil
}

func CreatePlugin(pv *models.PluginVO) (primitive.ObjectID, error) {
	p, err := convertToPluginDB(pv)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return db.CreatePlugin(p)
}

func CheckPluginFileStatus(id, language string) error {
	p, err := GetPluginByID(id)
	if err != nil {
		return err
	}
	return shared.NewPluginFactory(p.ID, p.Code, p.Language, *configs.PluginRootPath).CheckFile()
}

func GetPluginByID(id string) (plugin models.PluginVO, err error) {
	p, err := db.GetPluginByID(id)
	if err != nil {
		return
	}
	return convertToPluginVO(p), nil
}

func UpdatePlugin(pv *models.PluginVO) error {
	p, err := convertToPluginDB(pv)
	if err != nil {
		return err
	}
	return db.UpdatePlugin(p)
}

func DeletePlugin(id string) error {
	if err := db.DeletePlugin(id); err != nil {
		return err
	}
	return nil
}

func GetPluginPagination(q bindmodels.Query) (models.Pagination, error) {
	result := models.Pagination{
		Items: []models.PipelineVO{},
		Total: 0,
	}
	pvos, err := GetPlugins(q)
	if err != nil {
		return result, err
	}
	total, err := db.GetPluginCount()
	if err != nil {
		return result, err
	}
	result.Items = pvos
	result.Total = total
	return result, nil
}

func ReceiveBinaryFile(data []byte, id string) {
	// db.UpdatePluginFile(data, id)
}

func DebugPlugin(pv models.PluginVO) (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("Recovered in %s", r)
		}
	}()

	pf := shared.NewPluginFactory(pv.FileID, pv.Code, pv.Language, *configs.PluginRootPath)
	pc, err := pf.GetPluginClient()
	defer pc.Close()
	if err := pc.Check(); err != nil {
		result = err.Error()
		return
	}
	p, err := pc.AcquirePlugin("")
	defer pc.ReleasePlugin(p)
	if err != nil {
		result = err.Error()
		return
	}
	data := make(map[string]string, len(pv.Inputs))
	for _, input := range pv.Inputs {
		data[input.Name] = input.Value
	}
	if rr, err := p.Run(data); err != nil {
		result = err.Error()
		return
	} else {
		jsonString, err := json.Marshal(rr)
		if err != nil {
			result = err.Error()
			return
		}
		return string(jsonString)
	}
}

func PublishPlugin(pluginId string) error {
	p, err := db.GetPluginByID(pluginId)
	if err != nil {
		return err
	}
	if err := db.UpdatePlugin(&p); err != nil {
		return err
	}
	pf := shared.NewPluginFactory(p.FileID.Hex(), p.Code, p.Language, *configs.PluginRootPath)
	pc, err := pf.GetPluginClientFromDB(db.DownloadPluginFile)
	if err != nil {
		return err
	}
	err = shared.PM.UpdatePlugin(pluginId, pc)
	if err != nil {
		return err
	}
	return nil
}

func CheckPluginHealth(id string) string {
	return shared.PM.CheckPluginHealth(id)
}

func convertToPluginVO(p db.Plugin) models.PluginVO {
	pv := models.PluginVO{
		ID:         p.ID.Hex(),
		Name:       p.Name,
		Desc:       p.Desc,
		Language:   p.Language,
		Code:       p.Code,
		CreateTime: p.CreateTime.In(common.GetLocalZone()).Format("2006-01-02 15:04:05"),
		UpdateTime: p.UpdateTime.In(common.GetLocalZone()).Format("2006-01-02 15:04:05"),
	}
	var inputs []models.Input
	for _, input := range p.Inputs {
		inputs = append(inputs, models.Input{
			Name:     input.Name,
			Desc:     input.Name,
			Required: input.Required,
			Value:    input.Value,
		})
	}
	pv.Inputs = inputs
	return pv
}

func convertToPluginDB(pv *models.PluginVO) (*db.Plugin, error) {
	var inputs []db.Input
	for _, input := range pv.Inputs {
		inputs = append(inputs, db.Input{
			Name:     input.Name,
			Desc:     input.Desc,
			Required: input.Required,
			Value:    input.Value,
		})
	}
	p := &db.Plugin{
		Name:     pv.Name,
		Inputs:   inputs,
		Code:     pv.Code,
		Desc:     pv.Desc,
		Language: pv.Language,
	}
	if pv.ID == "" {
		return p, nil
	} else {
		oid, err := db.GetIDFromInterface(pv.ID)
		if err != nil {
			return nil, err
		}
		p.ID = oid
		return p, nil
	}
}
