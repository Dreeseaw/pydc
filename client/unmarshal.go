/*
    main code for unmarshalling
    python3's marshal.dumps()
*/
package main

import (
    "fmt"
    // "strconv"
    "encoding/binary"
    "math"
)

var (
    TUPLE_TYPE_BYTE  uint8 = 169
    INT_TYPE_BYTE    uint8 = 233
    FLOAT_TYPE_BYTE  uint8 = 231 
    STR_TYPE_BYTE    uint8 = 218
    TRUE_TYPE_BYTE   uint8 = 84
    FALSE_TYPE_BYTE  uint8 = 70
)

// cut and return 1 byte from p
func r_byte(p []byte, index *int) byte {
    ret := p[*index]
    *index += 1
    return ret
}

// cut and return str_len bytes from p
func r_string(p []byte, str_len int, index *int) []byte {
    ret := p[*index:*index+str_len]
    *index += str_len
    return ret
}

// cut and return int32 from p
func r_long(p []byte, index *int) int32 {
    long_bytes := r_string(p, 4, index)
    var ret int32 = -1
    ret = (int32)(long_bytes[0])
    ret = ret | ((int32)(long_bytes[1]) << 8)
    ret = ret | ((int32)(long_bytes[2]) << 16)
    ret = ret | ((int32)(long_bytes[3]) << 24)
    return ret
}

// cut and return float64 from p
func r_float(p []byte, index *int) float64 {
    float_bytes := r_string(p, 8, index)
    bits := binary.LittleEndian.Uint64(float_bytes)
    return math.Float64frombits(bits)
}

// given payload, return a mock python object
func r_object(payload []byte, index *int) ([]interface{}, error) {

    curIndex := 0
    tmp_object := make([]interface{}, 1)
        
    obj_type_byte := r_byte(payload, &curIndex)
    obj_type := (uint8)(obj_type_byte)

    if obj_type == TUPLE_TYPE_BYTE {
        tuple_len_byte := r_byte(payload, &curIndex)
        tuple_len := (uint8)(tuple_len_byte)
        if tuple_len == 0 {
            return nil, fmt.Errorf("tuple len is 0")
        }

        tmp_tuple := make([]interface{}, tuple_len)

        for tuple_pos, _ := range tmp_tuple {
            tmp_tuple_obj := r_object(payload, curIndex)
        }

    } else if obj_type == INT_TYPE_BYTE {
        val := r_long(payload, &curIndex)
        tmp_tuple[tuple_pos] = (int32)(val)

    } else if obj_type == FLOAT_TYPE_BYTE {
        val := r_float(payload, &curIndex)
        tmp_tuple[tuple_pos] = (float64)(val)
    
    } else if obj_type == STR_TYPE_BYTE {
        str_len := r_byte(payload, &curIndex)
        val := r_string(payload, (int)(str_len), &curIndex)
        tmp_tuple[tuple_pos] = (string)(val)
    
    } else if obj_type == TRUE_TYPE_BYTE {
        tmp_tuple[tuple_pos] = true

    } else if obj_type == FALSE_TYPE_BYTE {
        tmp_tuple[tuple_pos] = false
    
    } else {
        return nil, fmt.Errorf("type not recognized: %v", (uint8)(obj_type))
    }

    return tmp_tuple, nil
}
