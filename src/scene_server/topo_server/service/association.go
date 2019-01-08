/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

// CreateMainLineObject create a new object in the main line topo
func (s *topoService) CreateMainLineObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	mainLineAssociation := &metadata.Association{}

	_, err := mainLineAssociation.Parse(data)
	if nil != err {
		blog.Errorf("[api-asst] failed to parse the data(%#v), error info is %s", data, err.Error())
	}
	params.MetaData = &mainLineAssociation.Metadata
	return s.core.AssociationOperation().CreateMainlineAssociation(params, mainLineAssociation)
}

// DeleteMainLineObject delete a object int the main line topo
func (s *topoService) DeleteMainLineObject(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	objID := pathParams("bk_obj_id")

	params.MetaData = metadata.NewMetaDataFromInterface(data[metadata.BKMetadata])

	err := s.core.AssociationOperation().DeleteMainlineAssociaton(params, objID)
	return nil, err
}

// SearchMainLineOBjectTopo search the main line topo
func (s *topoService) SearchMainLineObjectTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	params.MetaData = metadata.NewMetaDataFromInterface(data[metadata.BKMetadata])

	bizObj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("[api-asst] failed to find the biz object, error info is %s", err.Error())
		return nil, err
	}

	return s.core.AssociationOperation().SearchMainlineAssociationTopo(params, bizObj)
}

// SearchObjectByClassificationID search the object by classification ID
func (s *topoService) SearchObjectByClassificationID(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	bizObj, err := s.core.ObjectOperation().FindSingleObject(params, pathParams("bk_obj_id"))
	if nil != err {
		blog.Errorf("[api-asst] failed to find the biz object, error info is %s", err.Error())
		return nil, err
	}

	params.MetaData = metadata.NewMetaDataFromInterface(data[metadata.BKMetadata])

	return s.core.AssociationOperation().SearchMainlineAssociationTopo(params, bizObj)
}

// SearchBusinessTopo search the business topo
func (s *topoService) SearchBusinessTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	paramPath := frtypes.MapStr{}
	paramPath.Set("id", pathParams("bk_biz_id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-asst] failed to parse the path params id(%s), error info is %s ", pathParams("app_id"), err.Error())
		return nil, err
	}

	bizObj, err := s.core.ObjectOperation().FindSingleObject(params, common.BKInnerObjIDApp)
	if nil != err {
		return nil, err
	}

	params.MetaData = metadata.NewMetaDataFromInterface(data[metadata.BKMetadata])
	data.Remove(metadata.BKMetadata)

	return s.core.AssociationOperation().SearchMainlineAssociationInstTopo(params, bizObj, id)
}

// SearchMainLineChildInstTopo search the child inst topo by a inst
func (s *topoService) SearchMainLineChildInstTopo(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	//{obj_id}/{app_id}/{inst_id}
	objID := pathParams("obj_id")
	bizID, err := strconv.ParseInt(pathParams("app_id"), 10, 64)
	if nil != err {
		return nil, params.Err.Errorf(common.CCErrCommParamsIsInvalid, "app_id")
	}

	// get the instance id of this object.
	instID, err := strconv.ParseInt(pathParams("inst_id"), 10, 64)
	if nil != err {
		return nil, params.Err.Errorf(common.CCErrCommParamsIsInvalid, "inst_id")
	}
	_ = bizID

	obj, err := s.core.ObjectOperation().FindSingleObject(params, objID)
	if nil != err {
		return nil, err
	}

	return s.core.AssociationOperation().SearchMainlineAssociationInstTopo(params, obj, instID)
}

func (s *topoService) SearchAssociationType(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	request := &metadata.SearchAssociationTypeRequest{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}
	ret, err := s.core.AssociationOperation().SearchType(params, request)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (s *topoService) SearchObjectAssoWithAssoKindList(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	ids := new(metadata.AssociationKindIDs)
	if err := data.MarshalJSONInto(ids); err != nil {
		return nil, params.Err.Error(common.CCErrCommParamsInvalid)
	}

	return s.core.AssociationOperation().SearchObjectAssoWithAssoKindList(params, ids.AsstIDs)
}

func (s *topoService) CreateAssociationType(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	request := &metadata.AssociationKind{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}
	ret, err := s.core.AssociationOperation().CreateType(params, request)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (s *topoService) UpdateAssociationType(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	request := &metadata.UpdateAssociationTypeRequest{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}
	asstTypeID, err := strconv.Atoi(pathParams("id"))
	if err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}
	ret, err := s.core.AssociationOperation().UpdateType(params, asstTypeID, request)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (s *topoService) DeleteAssociationType(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	asstTypeID, err := strconv.Atoi(pathParams("id"))
	if err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	ret, err := s.core.AssociationOperation().DeleteType(params, asstTypeID)
	if err != nil {
		return nil, err
	}

	if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (s *topoService) SearchAssociationInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	request := &metadata.SearchAssociationInstRequest{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	params.MetaData = metadata.NewMetaDataFromInterface(data[metadata.BKMetadata])
	data.Remove(metadata.BKMetadata)

	ret, err := s.core.AssociationOperation().SearchInst(params, request)
	if err != nil {
		return nil, err
	} else if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	} else {
		return ret.Data, nil

	}
}

func (s *topoService) CreateAssociationInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {
	request := &metadata.CreateAssociationInstRequest{}
	if err := data.MarshalJSONInto(request); err != nil {
		return nil, params.Err.New(common.CCErrCommParamsInvalid, err.Error())
	}

	params.MetaData = metadata.NewMetaDataFromInterface(data[metadata.BKMetadata])

	ret, err := s.core.AssociationOperation().CreateInst(params, request)
	if err != nil {
		return nil, err
	} else if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	} else {
		return ret.Data, nil
	}
}

func (s *topoService) DeleteAssociationInst(params types.ContextParams, pathParams, queryParams ParamsGetter, data frtypes.MapStr) (interface{}, error) {

	id, err := strconv.ParseInt(pathParams("association_id"), 10, 64)
	if err != nil {
		return nil, params.Err.Error(common.CCErrCommParamsIsInvalid)
	}

	//biz id in delete object association
	params.MetaData = metadata.NewMetaDataFromInterface(data[metadata.BKMetadata])
	data.Remove(metadata.BKMetadata)

	ret, err := s.core.AssociationOperation().DeleteInst(params, id)
	if err != nil {
		return nil, err
	} else if ret.Code != 0 {
		return nil, params.Err.New(ret.Code, ret.ErrMsg)
	} else {
		return ret.Data, nil

	}
}
