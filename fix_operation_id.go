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

	
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		fmt.Printf("ไม่สามารถอ่านไฟล์ได้: %v\n", err)
		return
	}

	content := string(data)


	updatedContent, changes := fixDuplicateOperationIDs(content)


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

	
	ext := filepath.Ext(yamlFile)
	baseName := strings.TrimSuffix(yamlFile, ext)
	newYamlFile := baseName + "_fixed" + ext

	
	err = ioutil.WriteFile(newYamlFile, []byte(updatedContent), 0644)
	if err != nil {
		fmt.Printf("ไม่สามารถเขียนไฟล์ได้: %v\n", err)
		return
	}

	fmt.Printf("แก้ไขเรียบร้อยแล้ว ไฟล์ใหม่: %s\n", newYamlFile)
}


func fixDuplicateOperationIDs(content string) (string, map[string]map[string]string) {
	
	lines := strings.Split(content, "\n")
	
	
	tagOperations := make(map[string]map[string]int)
	
	
	changes := make(map[string]map[string]string)
	
	
	currentTag := ""
	

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		
		if strings.Contains(line, "tags:") && i+1 < len(lines) && strings.Contains(lines[i+1], "-") {
			
			tagLine := lines[i+1]
			tagParts := strings.Split(tagLine, "-")
			if len(tagParts) > 1 {
				tagName := strings.TrimSpace(tagParts[1])
				currentTag = tagName
			}
		}
		
		
		if strings.Contains(line, "operationId:") {
			
			parts := strings.Split(line, "operationId:")
			if len(parts) > 1 {
				operationId := strings.TrimSpace(parts[1])
				
			
				tag := "default"
				if currentTag != "" {
					tag = currentTag
				}
				
			
				if _, exists := tagOperations[tag]; !exists {
					tagOperations[tag] = make(map[string]int)
				}
				
				if _, exists := changes[tag]; !exists {
					changes[tag] = make(map[string]string)
				}
				
				if count, exists := tagOperations[tag][operationId]; exists {
					newOperationId := fmt.Sprintf("%s%d", operationId, count)
					lines[i] = strings.Replace(line, operationId, newOperationId, 1)
					tagOperations[tag][operationId] = count + 1
					
					changes[tag][operationId] = newOperationId
				} else {
					tagOperations[tag][operationId] = 1
				}
			}
		}
	}
	
	
	return strings.Join(lines, "\n"), changes
}

func countChanges(changes map[string]map[string]string) int {
	count := 0
	for _, operations := range changes {
		count += len(operations)
	}
	return count
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}