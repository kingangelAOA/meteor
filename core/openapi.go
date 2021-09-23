package core

// var Methods = []string{
// 	http.MethodConnect,
// 	http.MethodDelete,
// 	http.MethodGet,
// 	http.MethodHead,
// 	http.MethodOptions,
// 	http.MethodPatch,
// 	http.MethodPost,
// 	http.MethodPut,
// 	http.MethodTrace,
// }

// type SwaggerExpand struct {
// 	Swagger *openapi3.Swagger `json:"swagger,omitempty" bson:"swagger,omitempty"`
// }

// func (se *SwaggerExpand) GetOpenApiVersion() string {
// 	return se.Swagger.OpenAPI
// }

// func (se *SwaggerExpand) GetApiTitle() string {
// 	return se.Swagger.Info.Title
// }

// func (se *SwaggerExpand) GetPersonOrganizationName() string {
// 	return se.Swagger.Info.Contact.Name
// }

// func (se *SwaggerExpand) GetPersonOrganizationEmail() string {
// 	return se.Swagger.Info.Contact.Email
// }

// func (se *SwaggerExpand) GetPersonOrganizationUrl() string {
// 	return se.Swagger.Info.Contact.URL
// }

// func (se *SwaggerExpand) GetApiVersion() string {
// 	return se.Swagger.Info.Version
// }

// func (se *SwaggerExpand) GetServersUrls() []string {
// 	var urls []string
// 	for _, s := range se.Swagger.Servers {
// 		urls = append(urls, s.URL)
// 	}
// 	return urls
// }

// func (se *SwaggerExpand) GetTags() []string {
// 	var tags []string
// 	for _, t := range se.Swagger.Tags {
// 		tags = append(tags, t.Name)
// 	}
// 	return tags
// }

// func (se *SwaggerExpand) GetInterfaceNum() int {
// 	num := 0
// 	for _, v := range se.Swagger.Paths {
// 		for _, method := range Methods {
// 			if v.GetOperation(method) != nil {
// 				num++
// 			}
// 		}

// 	}
// 	return num
// }
