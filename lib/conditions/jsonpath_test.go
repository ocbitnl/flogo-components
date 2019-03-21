/*
* Copyright © 2017. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package condition

import (
	"fmt"
	"testing"
)

const JSONCONTENT = `{"category":{"id":16,"name":"Animals"},"id":16,"name":"SPARROW","photoUrls":["string"],"status":"sold","tags":[{"id":0,"name":"string"}]}`

func TestJSONPath(t *testing.T) {
	res, err := JsonPathEval(JSONCONTENT, "$.name")
	fmt.Println("res", *res)
	fmt.Println("err", err)
}
