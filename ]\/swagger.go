package controllers

// "meteor/db"

//UpdateSwagger update swagger with project
// func UpdateSwagger(c *gin.Context) {
// 	var cusr models.CreateUpdateSwaggerRequest
// 	if err := c.BindJSON(&cusr); err != nil {
// 		models.ResponseErr(c, err.Error())
// 		return
// 	}
// 	// if err := db.CreateOrUpdateSwagger(cusr); err != nil {
// 	// 	models.ResponseErr(c, err.Error())
// 	// 	return
// 	// }
// 	models.ResponseSuccessNoData(c)
// }

// //GetSwaggerVersions get swagger versions by projectId
// func GetSwaggerVersions(c *gin.Context) {
// 	var qs models.QuerySwaggers
// 	if err := c.BindQuery(&qs); err != nil {
// 		models.ResponseErr(c, err.Error())
// 		return
// 	}
// 	oid, err := primitive.ObjectIDFromHex(qs.ProjectID)
// 	if err != nil {
// 		models.ResponseErr(c, err.Error())
// 		return
// 	}
// 	// versions, err := db.GetSwaggerVersions(oid)
// 	// if err != nil {
// 	// 	models.ResponseErr(c, err.Error())
// 	// 	return
// 	// }
// 	// models.ResponseSuccess(c, versions)
// }

// func GetSwaggerYaml(c *gin.Context) {
// 	// var qs models.QuerySwaggers
// 	// if err := c.BindQuery(&qs); err != nil {
// 	// 	models.ResponseErr(c, err.Error())
// 	// 	return
// 	// }
// 	// oid, err := primitive.ObjectIDFromHex(qs.ProjectID)
// 	// if err != nil {
// 	// 	models.ResponseErr(c, err.Error())
// 	// 	return
// 	// }
// 	// r, err := db.GetSwagger(oid, qs.Version)
// 	// if err != nil {
// 	// 	models.ResponseErr(c, err.Error())
// 	// 	return
// 	// }
// 	// models.ResponseSuccess(c, r)
// }
