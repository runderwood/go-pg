//
// pg.go
//
package pg

// #include "pg_wrapper.h"
import "C"

import (
    "unsafe";
)

type PgResults struct {
    result unsafe.Pointer;
    row    int;
}

func Connect(conninfo string) unsafe.Pointer {
    conn := C.PgConnectDb(C.CString(conninfo));
    return conn;
}

func Close(conn unsafe.Pointer) {
    C.PgFinish(conn);
}

func Status(conn unsafe.Pointer) int {
    status := C.PgStatus(conn);
    return int(status);
}

func Exec(conn unsafe.Pointer, command string) *PgResults {
    res := new(PgResults);
    res.result = C.PgExec(conn, C.CString(command));
    res.row    = 0;
    return res;
}

func GetResult(conn unsafe.Pointer) unsafe.Pointer {
    return C.PgGetResult(conn);
}

func NFields(res *PgResults) int {
    return int(C.PgNFields(res.result));
}

func NTuples(res *PgResults) int {
    return int(C.PgNTuples(res.result));
}

func GetIsNull(res *PgResults, row_number int, column_number int) bool {
    if 1 == int(C.PgGetIsNull(res.result, _C_int(row_number), 
                                          _C_int(column_number))) {
      return true
    }
    return false
}

func GetValue(res *PgResults, row_number int, column_number int) string {
    return C.GoString(C.PgGetValue(res.result, _C_int(row_number), 
                                               _C_int(column_number)));
}

func GetLength(res *PgResults, row_number int, column_number int) int {
    return int(C.PgGetLength(res.result, _C_int(row_number), 
                                         _C_int(column_number)));
}

func FName(res *PgResults, column_number int) string {
    return C.GoString(C.PgFName(res.result, _C_int(column_number)));
}

func FNumber(res *PgResults, column_name string) int {
    return int(C.PgFNumber(res.result, C.CString(column_name)));
}

func FType(res *PgResults, column_number int) int {
    return int(C.PgFType(res.result, _C_int(column_number)));
}

func FSize(res *PgResults, column_number int) int {
    return int(C.PgFSize(res.result, _C_int(column_number)));
}

func FMod(res *PgResults, column_number int) int {
    return int(C.PgFMod(res.result, _C_int(column_number)));
}

func FetchRow(res *PgResults) []string {
    cols := NFields(res);
    col := 0;
    vals := make([]string, cols);
    if NTuples(res) > res.row {
        for col < cols {
            if GetIsNull(res, res.row, col) {
                vals[col] = "";
            } else {
                vals[col] = GetValue(res, res.row, col);
            }
            col += 1;
        }
        res.row += 1;
        return vals;
    }

    return nil;
}

