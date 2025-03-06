package model

import "gorm.io/gorm"

// import "fmt"

type Dictionary struct {
	Id         int        `gorm:"column:id;primaryKey" json:"id"`
	Definition string     `gorm:"column:definition" json:"definition"`
	Categories []Category `gorm:"many2many:dictionaries_categories" json:"categories"`
	Image      string     `gorm:"column:image" json:"image"`
	AudioGB    string     `gorm:"column:audio_gb" json:"audio_gb"`
	AudioUS    string     `gorm:"column:audio_us" json:"audio_us"`

	gorm.Model
}

// CONSTRUCTOR
func NewDictionary(definition string, categories []Category, img string) *Dictionary {
	return &Dictionary{Definition: definition, Categories: categories, Image: img}
}

// // METHODS
// func (d *Dictionary) Save() {
// 	res := database.DB.Create(d)

// 	if res.Error != nil {
// 		fmt.Printf("❌ Save '%s' DICTIONARY failed.\n", d.Definition)
// 		return
// 	}
// 	fmt.Printf("✅ Save '%s' DICTIONARY successfully.\n", d.Definition)

// }

// func (d *Dictionary) SaveImageToLocal(imgSrc string) {
// 	if imgSrc == "" {
// 		return
// 	}

// 	fileName := d.Definition2FileName()

// 	_, err := utils.DownloadImage(imgSrc, fileName)
// 	utils.CheckError(err)
// }

// func (d *Dictionary) SaveAudioToLocal(lang string) {
// 	fileName := d.Definition2FileName()

// 	if fileName == "" {
// 		return
// 	}

// 	if lang == "" {
// 		fileName = fmt.Sprintf("%s_us", fileName)
// 	} else {
// 		fileName = fmt.Sprintf("%s_%s", fileName, lang)
// 	}

// 	// call tts api
// 	speechData, err := tts.TTS(tts.TTS_CLIENT, lang, d.Definition)
// 	utils.CheckError(err)

// 	// save to local
// 	utils.Bytes2Audio(fileName, speechData)
// }

// func (d Dictionary) DisplayDictionary() {
// 	fmt.Printf("Definition: %s\nImage: %s\nGB: %s\nUS: %s\n", d.Definition, d.Image, d.AudioGB, d.AudioUS)
// }

// func (d Dictionary) Definition2FileName() string {
// 	return strings.ReplaceAll(strings.ReplaceAll(d.Definition, " ", "_"), "-", "_")
// }
