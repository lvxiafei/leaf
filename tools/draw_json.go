package tools

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func (j *JSONData) WriteToJson(src string) {

	data, err := json.MarshalIndent(j, "", "  ") // 第二个表示每行的前缀，这里不用，第三个是缩进符号
	checkError(err)

	err = os.WriteFile(src, data, 0777)
	checkError(err)
}

func (j *JSONData) OpenToJson(src string, jsonData JSONData) {

	file, err := os.ReadFile(src)
	if err != nil {
		log.Fatalf("Some error occured while reading file. Error: %s", err)
	}
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}
	res, err := PrettyStruct(jsonData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("res Struct: %s\n", res)
}

func (j *JSONData) OpenToAppend(src string) {

	var extraJsonData JSONData
	file, err := os.ReadFile(src)
	if err != nil {
		log.Fatalf("Some error occured while reading file. Error: %s", err)
	}
	err = json.Unmarshal(file, &extraJsonData)
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}
	j.Elements = append(j.Elements, extraJsonData.Elements...)
}

func (j *JSONData) OpenSceneHookFsToAppend(fs embed.FS, src string) {

	var extraJsonData JSONData
	file, err := fs.ReadFile("scene-hook/" + src)
	if err != nil {
		log.Fatalf("Some error occured while reading file. Error: %s", err)
	}
	err = json.Unmarshal(file, &extraJsonData)
	if err != nil {
		log.Fatalf("Error occured during unmarshaling. Error: %s", err.Error())
	}
	j.Elements = append(j.Elements, extraJsonData.Elements...)
}

func (j *JSONData) AddDownArrowRectangleWithText(jsonData *JSONData, headRectangle *Elements, text string, index string) *Elements {

	nextRectangle := *headRectangle
	nextRectangle.ID = "rectangle" + "-" + index
	nextRectangle.Y = headRectangle.Y + headRectangle.Height + 70

	nextArrow := Elements{
		Type:            "arrow",
		Version:         0,
		VersionNonce:    0,
		IsDeleted:       false,
		ID:              "",
		FillStyle:       "hachure",
		StrokeWidth:     1,
		StrokeStyle:     "solid",
		Roughness:       1,
		Opacity:         100,
		Angle:           0,
		X:               0,
		Y:               70,
		StrokeColor:     "#1e1e1e",
		BackgroundColor: "transparent",
		Width:           2,
		Height:          70,
		Seed:            2053740462,
		GroupIds:        nil,
		FrameID:         nil,
		Roundness: Roundness{
			Type: 2,
		},
		BoundElements: nil,
		Updated:       0,
		Link:          nil,
		Locked:        false,
		FontSize:      0,
		FontFamily:    0,
		Text:          "",
		TextAlign:     "",
		VerticalAlign: "",
		ContainerID:   nil,
		OriginalText:  "",
		LineHeight:    0,
		Baseline:      0,
		StartBinding: StartBinding{
			ElementID: headRectangle.ID,
			Focus:     0.03,
			Gap:       1,
		},
		EndBinding: EndBinding{
			ElementID: nextRectangle.ID,
			Focus:     0.03,
			Gap:       10,
		},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       nil,
		Points: [][]float64{
			{0, 0},
			{1, 60},
		},
	}
	nextArrow.EndArrowhead = "arrow"
	nextArrow.X = headRectangle.Width / 2
	nextArrow.Y = headRectangle.Y + 70
	nextArrow.ID = "arrow" + "-" + index

	headRectangle.BoundElements = nil
	headRectangle.BoundElements = append(headRectangle.BoundElements, BoundElements{
		ID:   nextArrow.ID,
		Type: "arrow",
	})
	for i, element := range jsonData.Elements {
		if element.ID == headRectangle.ID {
			jsonData.Elements[i].BoundElements = append(jsonData.Elements[i].BoundElements, headRectangle.BoundElements...)
		}
	}

	nextRectangle.BoundElements = nil
	nextRectangle.BoundElements = append(nextRectangle.BoundElements, BoundElements{
		ID:   nextArrow.ID,
		Type: "arrow",
	})

	nextRectangleText := Elements{
		Type:               "text",
		Version:            0,
		VersionNonce:       0,
		IsDeleted:          false,
		ID:                 "",
		FillStyle:          "hachure",
		StrokeWidth:        1,
		StrokeStyle:        "solid",
		Roughness:          1,
		Opacity:            100,
		Angle:              0,
		X:                  0,
		Y:                  0,
		StrokeColor:        "#1e1e1e", // black
		BackgroundColor:    "transparent",
		Width:              60,
		Height:             24,
		Seed:               2053740462,
		GroupIds:           nil,
		FrameID:            nil,
		Roundness:          Roundness{},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           20,
		FontFamily:         3,
		Text:               "",
		TextAlign:          "center",
		VerticalAlign:      "middle",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         1.2,
		Baseline:           19,
		IsFrameName:        false,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       "",
		Points:             nil,
	}
	nextRectangleText.ID = "text" + "-" + index
	nextRectangleText.X = nextRectangle.X + 70
	nextRectangleText.Y = nextRectangle.Y + 23
	nextRectangleText.ContainerID = nextRectangle.ID
	nextRectangleText.Text, nextRectangleText.OriginalText = text, text

	headRectangle.BoundElements = nil
	nextRectangle.BoundElements = append(nextRectangle.BoundElements, BoundElements{
		ID:   nextRectangleText.ID,
		Type: "text",
	})

	jsonData.Elements = append(jsonData.Elements, nextRectangle)
	jsonData.Elements = append(jsonData.Elements, nextRectangleText)
	jsonData.Elements = append(jsonData.Elements, nextArrow)
	return &nextRectangle
}

func (j *JSONData) AddUpArrowRectangleWithTextRightX(jsonData *JSONData, headRectangle *Elements, text string, index string, rightX float64) *Elements {

	nextRectangle := *headRectangle
	nextRectangle.ID = "rectangleUp" + "-" + index
	nextRectangle.Y = headRectangle.Y - headRectangle.Height - 70

	nextArrow := Elements{
		Type:            "arrow",
		Version:         0,
		VersionNonce:    0,
		IsDeleted:       false,
		ID:              "",
		FillStyle:       "hachure",
		StrokeWidth:     1,
		StrokeStyle:     "solid",
		Roughness:       1,
		Opacity:         100,
		Angle:           0,
		X:               0,
		Y:               70,
		StrokeColor:     "#1e1e1e",
		BackgroundColor: "transparent",
		Width:           2,
		Height:          70,
		Seed:            2053740462,
		GroupIds:        nil,
		FrameID:         nil,
		Roundness: Roundness{
			Type: 2,
		},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           0,
		FontFamily:         0,
		Text:               "",
		TextAlign:          "",
		VerticalAlign:      "",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         0,
		Baseline:           0,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       nil,
		Points: [][]float64{
			{0, 0},
			{1, 60},
		},
	}
	nextArrow.StartBinding = StartBinding{
		ElementID: nextRectangle.ID,
		Focus:     0.03,
		Gap:       10,
	}
	nextArrow.EndBinding = EndBinding{
		ElementID: headRectangle.ID,
		Focus:     0.03,
		Gap:       1,
	}
	nextArrow.StartArrowhead = "arrow"
	nextArrow.X = headRectangle.Width/2 + rightX
	nextArrow.Y = headRectangle.Y - 70
	nextArrow.ID = "arrowUp" + "-" + index

	headRectangle.BoundElements = nil
	headRectangle.BoundElements = append(headRectangle.BoundElements, BoundElements{
		ID:   nextArrow.ID,
		Type: "arrow",
	})
	for i, element := range jsonData.Elements {
		if element.ID == headRectangle.ID {
			jsonData.Elements[i].BoundElements = append(jsonData.Elements[i].BoundElements, headRectangle.BoundElements...)
		}
	}

	nextRectangle.BoundElements = nil
	nextRectangle.BoundElements = append(nextRectangle.BoundElements, BoundElements{
		ID:   nextArrow.ID,
		Type: "arrow",
	})

	nextRectangleText := Elements{
		Type:               "text",
		Version:            0,
		VersionNonce:       0,
		IsDeleted:          false,
		ID:                 "",
		FillStyle:          "hachure",
		StrokeWidth:        1,
		StrokeStyle:        "solid",
		Roughness:          1,
		Opacity:            100,
		Angle:              0,
		X:                  0,
		Y:                  0,
		StrokeColor:        "#1e1e1e", // black
		BackgroundColor:    "transparent",
		Width:              60,
		Height:             24,
		Seed:               2053740462,
		GroupIds:           nil,
		FrameID:            nil,
		Roundness:          Roundness{},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           20,
		FontFamily:         3,
		Text:               "",
		TextAlign:          "center",
		VerticalAlign:      "middle",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         1.2,
		Baseline:           19,
		IsFrameName:        false,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       "",
		Points:             nil,
	}
	nextRectangleText.ID = "textUp" + "-" + index
	nextRectangleText.X = nextRectangle.X - 70
	nextRectangleText.Y = nextRectangle.Y - 23
	nextRectangleText.ContainerID = nextRectangle.ID
	nextRectangleText.Text, nextRectangleText.OriginalText = text, text

	headRectangle.BoundElements = nil
	nextRectangle.BoundElements = append(nextRectangle.BoundElements, BoundElements{
		ID:   nextRectangleText.ID,
		Type: "text",
	})

	jsonData.Elements = append(jsonData.Elements, nextRectangle)
	jsonData.Elements = append(jsonData.Elements, nextRectangleText)
	jsonData.Elements = append(jsonData.Elements, nextArrow)
	return &nextRectangle
}

func (j *JSONData) AddUpArrowRectangleWithText(jsonData *JSONData, headRectangle *Elements, text string, index string) *Elements {

	nextRectangle := *headRectangle
	nextRectangle.ID = "rectangleUp" + "-" + index
	nextRectangle.Y = headRectangle.Y - headRectangle.Height - 70

	nextArrow := Elements{
		Type:            "arrow",
		Version:         0,
		VersionNonce:    0,
		IsDeleted:       false,
		ID:              "",
		FillStyle:       "hachure",
		StrokeWidth:     1,
		StrokeStyle:     "solid",
		Roughness:       1,
		Opacity:         100,
		Angle:           0,
		X:               0,
		Y:               70,
		StrokeColor:     "#1e1e1e",
		BackgroundColor: "transparent",
		Width:           2,
		Height:          70,
		Seed:            2053740462,
		GroupIds:        nil,
		FrameID:         nil,
		Roundness: Roundness{
			Type: 2,
		},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           0,
		FontFamily:         0,
		Text:               "",
		TextAlign:          "",
		VerticalAlign:      "",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         0,
		Baseline:           0,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       nil,
		Points: [][]float64{
			{0, 0},
			{1, 60},
		},
	}
	nextArrow.StartBinding = StartBinding{
		ElementID: nextRectangle.ID,
		Focus:     0.03,
		Gap:       10,
	}
	nextArrow.EndBinding = EndBinding{
		ElementID: headRectangle.ID,
		Focus:     0.03,
		Gap:       1,
	}
	nextArrow.StartArrowhead = "arrow"
	nextArrow.X = headRectangle.Width / 2
	nextArrow.Y = headRectangle.Y - 70
	nextArrow.ID = "arrowUp" + "-" + index

	headRectangle.BoundElements = nil
	headRectangle.BoundElements = append(headRectangle.BoundElements, BoundElements{
		ID:   nextArrow.ID,
		Type: "arrow",
	})
	for i, element := range jsonData.Elements {
		if element.ID == headRectangle.ID {
			jsonData.Elements[i].BoundElements = append(jsonData.Elements[i].BoundElements, headRectangle.BoundElements...)
		}
	}

	nextRectangle.BoundElements = nil
	nextRectangle.BoundElements = append(nextRectangle.BoundElements, BoundElements{
		ID:   nextArrow.ID,
		Type: "arrow",
	})

	nextRectangleText := Elements{
		Type:               "text",
		Version:            0,
		VersionNonce:       0,
		IsDeleted:          false,
		ID:                 "",
		FillStyle:          "hachure",
		StrokeWidth:        1,
		StrokeStyle:        "solid",
		Roughness:          1,
		Opacity:            100,
		Angle:              0,
		X:                  0,
		Y:                  0,
		StrokeColor:        "#1e1e1e", // black
		BackgroundColor:    "transparent",
		Width:              60,
		Height:             24,
		Seed:               2053740462,
		GroupIds:           nil,
		FrameID:            nil,
		Roundness:          Roundness{},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           20,
		FontFamily:         3,
		Text:               "",
		TextAlign:          "center",
		VerticalAlign:      "middle",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         1.2,
		Baseline:           19,
		IsFrameName:        false,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       "",
		Points:             nil,
	}
	nextRectangleText.ID = "textUp" + "-" + index
	nextRectangleText.X = nextRectangle.X - 70
	nextRectangleText.Y = nextRectangle.Y - 23
	nextRectangleText.ContainerID = nextRectangle.ID
	nextRectangleText.Text, nextRectangleText.OriginalText = text, text

	headRectangle.BoundElements = nil
	nextRectangle.BoundElements = append(nextRectangle.BoundElements, BoundElements{
		ID:   nextRectangleText.ID,
		Type: "text",
	})

	jsonData.Elements = append(jsonData.Elements, nextRectangle)
	jsonData.Elements = append(jsonData.Elements, nextRectangleText)
	jsonData.Elements = append(jsonData.Elements, nextArrow)
	return &nextRectangle
}

func (j *JSONData) AddSingleText(text string) {
	count := strings.Count(text, "\n") + 1

	nextRectangleText := Elements{
		Type:               "text",
		Version:            0,
		VersionNonce:       0,
		IsDeleted:          false,
		ID:                 "",
		FillStyle:          "hachure",
		StrokeWidth:        1,
		StrokeStyle:        "solid",
		Roughness:          1,
		Opacity:            100,
		Angle:              0,
		X:                  0,
		Y:                  0,
		StrokeColor:        "#1e1e1e", // black
		BackgroundColor:    "transparent",
		Width:              2240,
		Height:             134.4,
		Seed:               2053740462,
		GroupIds:           nil,
		FrameID:            nil,
		Roundness:          Roundness{},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           20,
		FontFamily:         3,
		Text:               "",
		TextAlign:          "left",
		VerticalAlign:      "top",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         1.2,
		Baseline:           19,
		IsFrameName:        false,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       "",
		Points:             nil,
	}
	nextRectangleText.ID = "text" + "-" + "single"
	nextRectangleText.X = 1000
	nextRectangleText.Y = 0
	nextRectangleText.Width = 3340
	nextRectangleText.Height = 24 * float64(count)
	nextRectangleText.ContainerID = nil
	nextRectangleText.Text, nextRectangleText.OriginalText = text, text

	j.Elements = append(j.Elements, nextRectangleText)
}

func NewJsonData() *JSONData {
	return &JSONData{
		Type:    "excalidraw",
		Version: 0,
		Source:  "",
		AppState: AppState{
			GridSize:            nil,
			ViewBackgroundColor: "#ffffff",
		},
	}
}

func (j *JSONData) GenerateJsonFile(requestItems []string, replyItems []string, generateFile string) error {

	if len(requestItems) < 1 {
		return errors.New("len less than 1")
	}
	nextRectangle := j.AddHeadRectangle(requestItems[0])
	for i, item := range requestItems[1:] {
		nextRectangle = j.AddDownArrowRectangleWithText(j, nextRectangle, item, strconv.Itoa(i+1))
	}

	if len(replyItems) >= 1 {
		nextRectangle = j.AddHeadRectangleLastRightX(nextRectangle, replyItems[0], 400)
		for i, item := range replyItems[1:] {
			nextRectangle = j.AddUpArrowRectangleWithTextRightX(j, nextRectangle, item, strconv.Itoa(i+1), 400)
		}
	}
	j.WriteToJson(generateFile)

	return nil
}

func (j *JSONData) AddHeadRectangle(item string) *Elements {

	headRectangle := Elements{
		Type:            "rectangle",
		Version:         0,
		VersionNonce:    0,
		IsDeleted:       false,
		ID:              "rectangle-head",
		FillStyle:       "hachure",
		StrokeWidth:     1,
		StrokeStyle:     "solid",
		Roughness:       1,
		Opacity:         100,
		Angle:           0,
		X:               0,
		Y:               0,
		StrokeColor:     "#1e1e1e",
		BackgroundColor: "transparent",
		Width:           200,
		Height:          70,
		Seed:            2053740462,
		GroupIds:        nil,
		FrameID:         nil,
		Roundness: Roundness{
			Type: 3,
		},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           0,
		FontFamily:         0,
		Text:               "",
		TextAlign:          "",
		VerticalAlign:      "",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         0,
		Baseline:           0,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       "",
		Points:             nil,
	}
	headRectangleText := Elements{
		Type:               "text",
		Version:            0,
		VersionNonce:       0,
		IsDeleted:          false,
		ID:                 "",
		FillStyle:          "hachure",
		StrokeWidth:        1,
		StrokeStyle:        "solid",
		Roughness:          1,
		Opacity:            100,
		Angle:              0,
		X:                  0,
		Y:                  0,
		StrokeColor:        "#1e1e1e",
		BackgroundColor:    "transparent",
		Width:              60,
		Height:             24,
		Seed:               2053740462,
		GroupIds:           nil,
		FrameID:            nil,
		Roundness:          Roundness{},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           20,
		FontFamily:         3,
		Text:               "",
		TextAlign:          "center",
		VerticalAlign:      "middle",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         1.2,
		Baseline:           19,
		IsFrameName:        false,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       "",
		Points:             nil,
	}
	headRectangleText.ID = "text-head"
	headRectangleText.X = headRectangle.X + 70
	headRectangleText.Y = headRectangle.Y + 23
	headRectangleText.ContainerID = headRectangle.ID
	headRectangleText.Text, headRectangleText.OriginalText = item, item
	headRectangle.BoundElements = append(headRectangle.BoundElements, BoundElements{
		ID:   headRectangleText.ID,
		Type: "text",
	})

	j.Elements = append(j.Elements, headRectangle)
	j.Elements = append(j.Elements, headRectangleText)

	//j.WriteToJson(generateFile)
	// TODO http json response

	return &headRectangle
}

func (j *JSONData) AddHeadRectangleRightX(item string, rightX float64) *Elements {

	headRectangle := Elements{
		Type:            "rectangle",
		Version:         0,
		VersionNonce:    0,
		IsDeleted:       false,
		ID:              "",
		FillStyle:       "hachure",
		StrokeWidth:     1,
		StrokeStyle:     "solid",
		Roughness:       1,
		Opacity:         100,
		Angle:           0,
		X:               0,
		Y:               0,
		StrokeColor:     "#1e1e1e",
		BackgroundColor: "transparent",
		Width:           200,
		Height:          70,
		Seed:            2053740462,
		GroupIds:        nil,
		FrameID:         nil,
		Roundness: Roundness{
			Type: 3,
		},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           0,
		FontFamily:         0,
		Text:               "",
		TextAlign:          "",
		VerticalAlign:      "",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         0,
		Baseline:           0,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       "",
		Points:             nil,
	}
	headRectangle.ID = "rectangle-head-" + strconv.Itoa(int(rightX))
	headRectangle.X += rightX
	headRectangleText := Elements{
		Type:               "text",
		Version:            0,
		VersionNonce:       0,
		IsDeleted:          false,
		ID:                 "",
		FillStyle:          "hachure",
		StrokeWidth:        1,
		StrokeStyle:        "solid",
		Roughness:          1,
		Opacity:            100,
		Angle:              0,
		X:                  0,
		Y:                  0,
		StrokeColor:        "#1e1e1e",
		BackgroundColor:    "transparent",
		Width:              60,
		Height:             24,
		Seed:               2053740462,
		GroupIds:           nil,
		FrameID:            nil,
		Roundness:          Roundness{},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           20,
		FontFamily:         3,
		Text:               "",
		TextAlign:          "center",
		VerticalAlign:      "middle",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         1.2,
		Baseline:           19,
		IsFrameName:        false,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       "",
		Points:             nil,
	}

	headRectangleText.ID = "text-head-" + strconv.Itoa(int(rightX))
	headRectangleText.X = headRectangle.X + 70 + rightX
	headRectangleText.Y = headRectangle.Y + 23
	headRectangleText.ContainerID = headRectangle.ID
	headRectangleText.Text, headRectangleText.OriginalText = item, item
	headRectangle.BoundElements = append(headRectangle.BoundElements, BoundElements{
		ID:   headRectangleText.ID,
		Type: "text",
	})

	j.Elements = append(j.Elements, headRectangle)
	j.Elements = append(j.Elements, headRectangleText)

	return &headRectangle
}

func (j *JSONData) AddHeadRectangleLastRightX(nextRectangle *Elements, item string, rightX float64) *Elements {

	headRectangle := *nextRectangle
	headRectangle.ID = "rectangle-head-" + strconv.Itoa(int(rightX))
	headRectangle.X += rightX
	headRectangle.BoundElements = nil
	headRectangleText := Elements{
		Type:               "text",
		Version:            0,
		VersionNonce:       0,
		IsDeleted:          false,
		ID:                 "",
		FillStyle:          "hachure",
		StrokeWidth:        1,
		StrokeStyle:        "solid",
		Roughness:          1,
		Opacity:            100,
		Angle:              0,
		X:                  0,
		Y:                  0,
		StrokeColor:        "#1e1e1e",
		BackgroundColor:    "transparent",
		Width:              60,
		Height:             24,
		Seed:               2053740462,
		GroupIds:           nil,
		FrameID:            nil,
		Roundness:          Roundness{},
		BoundElements:      nil,
		Updated:            0,
		Link:               nil,
		Locked:             false,
		FontSize:           20,
		FontFamily:         3,
		Text:               "",
		TextAlign:          "center",
		VerticalAlign:      "middle",
		ContainerID:        nil,
		OriginalText:       "",
		LineHeight:         1.2,
		Baseline:           19,
		IsFrameName:        false,
		StartBinding:       StartBinding{},
		EndBinding:         EndBinding{},
		LastCommittedPoint: nil,
		StartArrowhead:     nil,
		EndArrowhead:       "",
		Points:             nil,
	}

	headRectangleText.ID = "text-head-" + strconv.Itoa(int(rightX))
	headRectangleText.X = headRectangle.X + 70 + rightX
	headRectangleText.Y = headRectangle.Y + 23
	headRectangleText.ContainerID = headRectangle.ID
	headRectangleText.Text, headRectangleText.OriginalText = item, item
	headRectangle.BoundElements = append(headRectangle.BoundElements, BoundElements{
		ID:   headRectangleText.ID,
		Type: "text",
	})

	j.Elements = append(j.Elements, headRectangle)
	j.Elements = append(j.Elements, headRectangleText)

	return &headRectangle
}
