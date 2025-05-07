package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("กรุณาระบุไฟล์ OpenAPI YAML: go run fix_operation_id.go <path-to-yaml-file>")
		return
	}

	yamlFile := os.Args[1]
	if !fileExists(yamlFile) {
		fmt.Printf("ไม่พบไฟล์: %s\n", yamlFile)
		return
	}

	// อ่านไฟล์ YAML เป็นข้อความ
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("ไม่สามารถอ่านไฟล์ได้: %v\n", err)
		return
	}

	// แปลงเป็นข้อความ
	content := string(data)

	// ตรวจสอบและแก้ไข operationId ที่ซ้ำกัน
	updatedContent, changes := fixDuplicateOperationIDs(content)

	// แสดงผลการแก้ไข
	if len(changes) > 0 {
		fmt.Println("\nรายการแก้ไข operationId ที่ซ้ำกัน:")
		fmt.Println("------------------------------------")
		for tag, operations := range changes {
			fmt.Printf("Tag: %s\n", tag)
			for original, modified := range operations {
				fmt.Printf("  - %s -> %s\n", original, modified)
			}
			fmt.Println()
		}
		fmt.Printf("จำนวนการแก้ไขทั้งหมด: %d รายการ\n\n", countChanges(changes))
	} else {
		fmt.Println("ไม่พบ operationId ที่ซ้ำกัน")
	}

	// สร้างชื่อไฟล์ใหม่
	ext := filepath.Ext(yamlFile)
	baseName := strings.TrimSuffix(yamlFile, ext)
	newYamlFile := baseName + "_fixed" + ext

	// เขียนข้อมูลลงไฟล์ใหม่
	err = ioutil.WriteFile(newYamlFile, []byte(updatedContent), 0644)
	if err != nil {
		fmt.Printf("ไม่สามารถเขียนไฟล์ได้: %v\n", err)
		return
	}

	fmt.Printf("แก้ไขเรียบร้อยแล้ว ไฟล์ใหม่: %s\n", newYamlFile)
}

// fixDuplicateOperationIDs แก้ไข operationId ที่ซ้ำกัน และคืนค่าการเปลี่ยนแปลง
func fixDuplicateOperationIDs(content string) (string, map[string]map[string]string) {
	// แยกบรรทัด
	lines := strings.Split(content, "\n")
	
	// เก็บ operationId ที่เจอแล้วตาม tag
	tagOperations := make(map[string]map[string]int)
	
	// เก็บการเปลี่ยนแปลง
	changes := make(map[string]map[string]string)
	
	// ตัวแปรสำหรับเก็บ tag ปัจจุบัน
	currentTag := ""
	
	// วนลูปผ่านทุกบรรทัด
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		
		// ตรวจสอบว่าเป็นบรรทัดที่มี tag หรือไม่
		if strings.Contains(line, "tags:") && i+1 < len(lines) && strings.Contains(lines[i+1], "-") {
			// อ่าน tag แรกในรายการ
			tagLine := lines[i+1]
			tagParts := strings.Split(tagLine, "-")
			if len(tagParts) > 1 {
				tagName := strings.TrimSpace(tagParts[1])
				currentTag = tagName
			}
		}
		
		// ตรวจสอบว่าเป็นบรรทัดที่มี operationId หรือไม่
		if strings.Contains(line, "operationId:") {
			// แยก operationId
			parts := strings.Split(line, "operationId:")
			if len(parts) > 1 {
				operationId := strings.TrimSpace(parts[1])
				
				// ถ้าไม่มี tag ให้ใช้ "default"
				tag := "default"
				if currentTag != "" {
					tag = currentTag
				}
				
				// สร้าง map สำหรับ tag ถ้ายังไม่มี
				if _, exists := tagOperations[tag]; !exists {
					tagOperations[tag] = make(map[string]int)
				}
				
				// สร้าง map สำหรับเก็บการเปลี่ยนแปลงถ้ายังไม่มี
				if _, exists := changes[tag]; !exists {
					changes[tag] = make(map[string]string)
				}
				
				// ตรวจสอบว่า operationId นี้เคยเจอใน tag นี้หรือยัง
				if count, exists := tagOperations[tag][operationId]; exists {
					// ถ้าเคยเจอแล้ว ให้เพิ่มตัวเลขต่อท้าย
					newOperationId := fmt.Sprintf("%s%d", operationId, count)
					lines[i] = strings.Replace(line, operationId, newOperationId, 1)
					tagOperations[tag][operationId] = count + 1
					
					// เก็บการเปลี่ยนแปลง
					changes[tag][operationId] = newOperationId
				} else {
					// ถ้ายังไม่เคยเจอ ให้เพิ่มเข้าไปใน map
					tagOperations[tag][operationId] = 1
				}
			}
		}
	}
	
	// รวมบรรทัดกลับเป็นข้อความ
	return strings.Join(lines, "\n"), changes
}

// countChanges นับจำนวนการเปลี่ยนแปลงทั้งหมด
func countChanges(changes map[string]map[string]string) int {
	count := 0
	for _, operations := range changes {
		count += len(operations)
	}
	return count
}

// fileExists ตรวจสอบว่าไฟล์มีอยู่จริงหรือไม่
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}