package gendoc

import "strconv"

var samplesCache = map[string]interface{}{}
var scalarTypes map[string]struct{}

func init() {
	scalars := makeScalars()
	scalarTypes = map[string]struct{}{}
	for _, scalar := range scalars {
		scalarTypes[scalar.ProtoType] = struct{}{}
	}
}

func SampleFilter(message *Message) string {
	panic("implement me!")
}

func generateSample(message *Message, file *File) interface{} {
	if cached, found := samplesCache[message.FullName]; found {
		return cached
	}

	if !message.HasFields {
		return struct{}{}
	}

	result := map[string]interface{}{}
	for _, field := range message.Fields {
		setMapField(result, field, generateFieldValue(field.FullType, field.IsMap, file))
	}
	samplesCache[message.FullName] = result
	return result
}

func generateFieldValue(typeName string, isMap bool, file *File) interface{} {
	if isScalarType(typeName) {
		return generateSampleScalarValue(typeName)
	}
	if isMap {
		return generateSampleMap(typeName, file)
	}
	fieldMessageIdx, found := file.MessageIndexes[typeName]
	if found {
		fieldMessage := file.Messages[fieldMessageIdx]
		return generateSample(fieldMessage, file)
	}
	fieldEnumIdx, found := file.EnumIndexes[typeName]
	if found {
		fieldEnum := file.Enums[fieldEnumIdx]
		return generateEnum(fieldEnum)
	}
	return nil
}

func setMapField(m map[string]interface{}, field *MessageField, value interface{}) {
	if field.Label == "repeated" {
		m[field.Name] = []interface{}{value}
		return
	}
	m[field.Name] = value
}

func generateEnum(enum *Enum) interface{} {
	number := enum.Values[len(enum.Values)-1].Number
	integer, _ := strconv.Atoi(number)
	return integer
}

func generateSampleMap(typeName string, file *File) interface{} {
	message := file.Messages[file.MessageIndexes[typeName]]
	result := map[interface{}]interface{}{}
	var keyType, valueType string
	for _, field := range message.Fields {
		if field.Name == "key" {
			keyType = field.FullType
		}
		if field.Name == "value" {
			valueType = field.FullType
		}
	}
	result[generateSampleScalarValue(keyType)] = generateFieldValue(valueType, false, file)

	return result
}

func generateSampleScalarValue(typeName string) interface{} {
	switch typeName {
	case "double", "float":
		return 0.01
	case "int32", "uint32", "sint32", "fixed32", "sfixed32":
		return 0
	case "int64", "uint64", "sint64", "fixed64", "sfixed64":
		return "0"
	case "bool":
		return false
	case "string":
		return "string"
	case "bytes":
		return "<binary bytes>"
	default:
		return "<unknown scalar value>"
	}
}

func isScalarType(typeName string) bool {
	_, is := scalarTypes[typeName]
	return is
}
