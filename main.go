package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Person struct {
	ID    string // RG ou CPF
	Name  string
	Phone string
}

type RegistryEntry struct {
	PersonID     string
	TimestampIn  time.Time
	TimestampOut time.Time
}

var peopleDB = make(map[string]Person)
var registryLog = make([]RegistryEntry, 0)
var activeEntries = make(map[string]int)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("SIGECO")
	myWindow.Resize(fyne.NewSize(700, 500))

	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("RG ou CPF (Obrigatório)")
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Nome Completo")
	phoneEntry := widget.NewEntry()
	phoneEntry.SetPlaceHolder("Telefone")

	statusLabel := widget.NewLabel("Aguardando ação...")
	statusLabel.Wrapping = fyne.TextWrapWord

	currentlyInsideData := binding.NewStringList()
	currentlyInsideList := widget.NewListWithData(
		currentlyInsideData,
		func() fyne.CanvasObject {
			return widget.NewLabel("template item")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)

	entryButton := widget.NewButton("Registrar Entrada", func() {
		id := idEntry.Text
		name := nameEntry.Text

		if id == "" {
			statusLabel.SetText("Erro: O ID (RG/CPF) é obrigatório.")
			return
		}

		if _, ok := activeEntries[id]; ok {
			statusLabel.SetText(fmt.Sprintf("Erro: %s (%s) já está dentro!", name, id))
			return
		}

		peopleDB[id] = Person{ID: id, Name: name, Phone: phoneEntry.Text}

		newEntry := RegistryEntry{
			PersonID:    id,
			TimestampIn: time.Now(),
		}

		registryLog = append(registryLog, newEntry)
		newEntryIndex := len(registryLog) - 1

		activeEntries[id] = newEntryIndex

		updateInsideListUI(currentlyInsideData)

		statusLabel.SetText(fmt.Sprintf("Entrada registrada: %s (%s)", name, id))
		clearFields(idEntry, nameEntry, phoneEntry)
	})

	exitButton := widget.NewButton("Registrar Saída", func() {
		id := idEntry.Text

		if id == "" {
			statusLabel.SetText("Erro: O ID (RG/CPF) é obrigatório.")
			return
		}

		entryIndex, ok := activeEntries[id]
		if !ok {
			statusLabel.SetText(fmt.Sprintf("Erro: Pessoa com ID %s não está registrada como 'dentro'.", id))
			return
		}

		registryLog[entryIndex].TimestampOut = time.Now()

		delete(activeEntries, id)

		updateInsideListUI(currentlyInsideData)

		personName := peopleDB[id].Name
		statusLabel.SetText(fmt.Sprintf("Saída registrada: %s (%s)", personName, id))
		clearFields(idEntry, nameEntry, phoneEntry)
	})

	form := widget.NewForm(
		widget.NewFormItem("ID (RG/CPF)", idEntry),
		widget.NewFormItem("Nome", nameEntry),
		widget.NewFormItem("Telefone", phoneEntry),
	)

	buttons := container.NewGridWithColumns(2, entryButton, exitButton)

	leftSide := container.NewVBox(
		form,
		buttons,
		statusLabel,
	)

	rightSide := container.NewBorder(
		widget.NewLabelWithStyle("Log de Eventos:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		currentlyInsideList,
	)

	split := container.NewHSplit(leftSide, rightSide)
	split.SetOffset(0.5)

	myWindow.SetContent(split)
	myWindow.ShowAndRun()
}
func updateInsideListUI(list binding.StringList) {
	var items []string
	
	for _, entry := range registryLog {
		personName := peopleDB[entry.PersonID].Name
		
		var itemString string

		if entry.TimestampOut.IsZero() {
			itemString = fmt.Sprintf("%s (%s) - Entrou: %s",
				personName,
				entry.PersonID,
				entry.TimestampIn.Format("15:04:05"),
			)
		} else {
			itemString = fmt.Sprintf("%s (%s) - SAIU: %s",
				personName,
				entry.PersonID,
				entry.TimestampOut.Format("15:04:05"),
			)
		}
		
		items = append(items, itemString)
	}
	
	list.Set(items)
}

func clearFields(entries ...*widget.Entry) {
	for _, entry := range entries {
		entry.SetText("")
	}
}
