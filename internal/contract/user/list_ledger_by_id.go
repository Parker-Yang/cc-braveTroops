package user

import (
	"encoding/json"
	"fmt"

	"github.com/Evolt0/cc-braveTroops/pkg"
	"github.com/Evolt0/def-braveTroops/consts/status"
	"github.com/Evolt0/def-braveTroops/proto"
	"github.com/Evolt0/def-braveTroops/proto/epkg"
	"github.com/Evolt0/def-braveTroops/proto/prefix"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func ListLedgerByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	data := &proto.BodyData{}
	err := pkg.Decode(args, data)
	if err != nil {
		return nil, err
	}
	list, err := execListLedgerByID(stub, data)
	if err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return nil, epkg.Wrapf(status.InternalServerError, "failed to marshal request: %v", err)
	}
	return marshal, nil
}

func execListLedgerByID(stub shim.ChaincodeStubInterface, data *proto.BodyData) (*proto.Ledger, error) {
	queryString := fmt.Sprintf("{\"selector\":{\"objectType\":\"%s\",\"sID\":\"%s\"}}", prefix.Amounts, data.ID)
	ledger := &proto.Ledger{}
	ledger.List = make([]proto.Amounts, 0)
	// 通过stub.GetQueryResult方法获取迭代器iterator
	iterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	// 延迟关闭迭代器iterator
	defer iterator.Close()
	for iterator.HasNext() {
		// 通过迭代器的Next()方法获取下一个对象的Key与Value值(*queryresult.KV)
		result, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		amounts := &proto.Amounts{}
		err = json.Unmarshal(result.Value, amounts)
		if err != nil {
			return nil, epkg.Wrapf(status.InternalServerError, "failed to unmarshal request: %v", err)
		}
		ledger.List = append(ledger.List, *amounts)
	}
	return ledger, nil
}
