package handler

import (
	"kiteq/binding"
	// "log"
	"sort"
)

//----------------持久化的handler
type DeliverPreHandler struct {
	BaseForwardHandler
	IEventProcessor
	exchanger *binding.BindExchanger
}

//------创建deliverpre
func NewDeliverPreHandler(name string, exchanger *binding.BindExchanger) *DeliverPreHandler {
	phandler := &DeliverPreHandler{}
	phandler.name = name
	phandler.exchanger = exchanger
	phandler.processor = phandler
	return phandler
}

func (self *DeliverPreHandler) TypeAssert(event IEvent) bool {
	_, ok := self.cast(event)
	return ok
}

func (self *DeliverPreHandler) cast(event IEvent) (val *DeliverEvent, ok bool) {
	val, ok = event.(*DeliverEvent)
	return
}

func (self *DeliverPreHandler) Process(ctx *DefaultPipelineContext, event IEvent) error {

	// log.Printf("DeliverPreHandler|Process|%t\n", event)

	pevent, ok := self.cast(event)
	if !ok {
		return ERROR_INVALID_EVENT_TYPE
	}

	binds := self.exchanger.FindBinds(pevent.Topic, pevent.MessageType, func(b *binding.Binding) bool {
		//过滤掉已经投递成功的分组
		return false
	})

	groupIds := make([]string, 0, len(binds))
	//按groupid归并
	for _, bind := range binds {
		//认为不存在，则加入到投递分组中 SearchString如果没找到是返回slice长度而不是-1
		if sort.SearchStrings(groupIds, bind.GroupId) == len(groupIds) {
			groupIds = append(groupIds, bind.GroupId)
		}
	}

	//如果没有可用的分组则直接跳过
	if len(groupIds) <= 0 {
		return nil
	}

	//过滤掉已经投递成功分组id
	pevent.DeliverGroups = groupIds

	ctx.SendForward(pevent)

	return nil

}
