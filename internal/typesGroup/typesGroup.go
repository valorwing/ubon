package typesGroup

type UBON_Type_Group uint8 // 2 bits (up to 4 groups)

const (
	UBON_Type_Group_Primitive                UBON_Type_Group = 0 // bool, int, float, string, etc.
	UBON_Type_Group_Object                   UBON_Type_Group = 1 // map/struct with key-value fields
	UBON_Type_Group_Array                    UBON_Type_Group = 2 // homogeneous arrays with all per-item fields required
	UBON_Type_Group_Array_With_Optional_Mask UBON_Type_Group = 3 // array of objects with per-item field masks
)
