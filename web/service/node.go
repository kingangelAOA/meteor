package service

import (
	"common/bindmodels"
	"web/db"
	"web/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetNodes(q bindmodels.Query) ([]models.NodeVO, error) {
	var nvs []models.NodeVO
	ns, err := db.GetNodes(q)
	if err != nil {
		return nvs, err
	}
	for i := 0; i < len(ns); i++ {
		nv := ns[i].GetNodeVO()
		nv.Index = i + 1 + int(q.Offset)
		nvs = append(nvs, nv)
	}
	return nvs, nil
}

func GetBindingPluginNodes(q bindmodels.Query) ([]models.NodeVO, error) {
	var nvs []models.NodeVO
	ns, err := db.GetBindingPluginNodes(q)
	if err != nil {
		return nvs, err
	}
	for i := 0; i < len(ns); i++ {
		nv := ns[i].GetNodeVO()
		nv.Index = i + 1 + int(q.Offset)

		nvs = append(nvs, nv)
	}
	return nvs, nil
}

func GetAllNodes() ([]models.NodeVO, error) {
	ns, err := db.GetAllNodes()
	if err != nil {
		return []models.NodeVO{}, err
	}
	return nsToNvs(ns)
}

func GetNodePagination(q bindmodels.Query) (models.Pagination, error) {
	result := models.Pagination{
		Items: []models.PipelineVO{},
		Total: 0,
	}
	pvos, err := GetNodes(q)
	if err != nil {
		return result, err
	}
	total, err := db.GetNodeCount()
	if err != nil {
		return result, err
	}
	result.Items = pvos
	result.Total = total
	return result, nil
}

func CreateNode(nv *models.NodeVO) (primitive.ObjectID, error) {
	n, err := convertToNodeDB(nv)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return db.CreateNode(n)
}

func UpdateNode(nv *models.NodeVO) error {
	n, err := convertToNodeDB(nv)
	if err != nil {
		return err
	}
	return db.UpdateNode(n)

}

func DeleteNode(ID string) error {

	if err := db.DeleteNode(ID); err != nil {
		return err
	}
	return nil
}

func convertToNodeVO(n db.Node) models.NodeVO {
	return models.NodeVO{
		ID:       n.ID.Hex(),
		Describe: n.Describe,
		Type:     n.Type,
	}
}

func convertToNodeDB(nv *models.NodeVO) (*db.Node, error) {
	pluginOid, err := db.GetIDFromInterface(nv.PluginID)
	if err != nil {
		return nil, err
	}
	n := &db.Node{
		Name:     nv.Name,
		Describe: nv.Describe,
		PluginID: pluginOid,
		Type:     nv.Type,
	}
	if nv.ID == "" {
		return n, nil
	} else {
		oid, err := db.GetIDFromInterface(nv.ID)
		if err != nil {
			return nil, err
		}
		n.ID = oid
		return n, nil
	}
}

func nsToNvs(ns []db.Node) ([]models.NodeVO, error) {
	var nvs []models.NodeVO
	for i := 0; i < len(ns); i++ {
		nv := ns[i].GetNodeVO()
		nv.Index = i + 1
		nvs = append(nvs, nv)
	}
	return nvs, nil
}
