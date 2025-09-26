package parser

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
)

func IsDex(addr *common.Address)(bool string){
	data,_:=ReadFile("config/DexAddr.json")

	found:=false

	var DeFiconfig map[string]interface{}
	json.Unmarshal(data, &DeFiconfig)

	dexMap := DeFiconfig["Dex"].(map[string]interface{})

	for name,dex :=range dexMap{
		switch v:=dex.type {
		case map[string]interface{}:


			for _ , detaildex := range v{

				if _,dexmethods=detaildex.(map[string]interface{});ok{

					for _,method :=range dexmethods{

						if method.(string)==addr{
							found:=true
							return found,name
						}
					}

				}else{

					if  s,ok:=detaildex.(string){
						found:=true
						return found,name
					}
					
				}
			}
			
		}
	}

	return found,name

}

func IsDeFiRouter(to *common.Address) (bool, string) {
	
	return IsDex(to)
}


