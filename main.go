package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Person struct {
	ID    string
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

type AppData struct {
	People   map[string]Person `json:"people"`
	Registry []RegistryEntry   `json:"registry"`
}

const dataFile = "sigeco_data.json"

const (
	FilterCompleto   = iota
	FilterDentro     = iota
	FilterTodos      = iota
	FilterUltimaHora = iota
	FilterDia        = iota
	FilterSaidas	 = iota
)

var currentFilterMode = FilterCompleto

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("SIGECO")

	loadData()

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

		saveData()
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

		saveData()
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

	filterDentroBtn := widget.NewButton("Quem está Dentro", func() {
		currentFilterMode = FilterDentro
		updateInsideListUI(currentlyInsideData)
	})
	filterTodosBtn := widget.NewButton("Visitantes", func() {
		currentFilterMode = FilterTodos
		updateInsideListUI(currentlyInsideData)
	})
	filterHoraBtn := widget.NewButton("Última Hora", func() {
		currentFilterMode = FilterUltimaHora
		updateInsideListUI(currentlyInsideData)
	})
	filterDiaBtn := widget.NewButton("Hoje", func() {
		currentFilterMode = FilterDia
		updateInsideListUI(currentlyInsideData)
	})
	filterCompletoBtn := widget.NewButton("Relatório Completo", func() {
		currentFilterMode = FilterCompleto
		updateInsideListUI(currentlyInsideData)
	})
	filterSaidasBtn := widget.NewButton("Saídas Realizadas", func() {
		currentFilterMode = FilterSaidas
		updateInsideListUI(currentlyInsideData)
	})

	filterButtons := container.NewGridWithColumns(3,
		filterDentroBtn,
		filterHoraBtn,
		filterDiaBtn,
		filterTodosBtn,
		filterCompletoBtn,
		filterSaidasBtn,
	)

	logTitle := widget.NewLabelWithStyle("Log de Eventos:", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	topContent := container.NewVBox(logTitle, filterButtons)

	rightSide := container.NewBorder(
		topContent,
		nil, nil, nil,
		currentlyInsideList,
	)

	split := container.NewHSplit(leftSide, rightSide)
	split.SetOffset(0.5)

	myWindow.SetContent(split)
	myWindow.ShowAndRun()
}

func formatLogEntry(entry RegistryEntry, layout string) string {
	personName := peopleDB[entry.PersonID].Name

	if entry.TimestampOut.IsZero() {
		return fmt.Sprintf("%s (%s) - Entrou: %s",
			personName,
			entry.PersonID,
			entry.TimestampIn.Format(layout),
		)
	} else {
		return fmt.Sprintf("%s (%s) - Entrou: %s | Saiu: %s",
			personName,
			entry.PersonID,
			entry.TimestampIn.Format(layout),
			entry.TimestampOut.Format(layout),
		)
	}
}

func updateInsideListUI(list binding.StringList) {
	var items []string
	const layout = "02/01/2006 15:04:05"
	now := time.Now()

	switch currentFilterMode {

	case FilterDentro:
		for id, entryIndex := range activeEntries {
			personName := peopleDB[id].Name
			timestamp := registryLog[entryIndex].TimestampIn
			itemString := fmt.Sprintf("%s (%s) - Entrou: %s",
				personName,
				id,
				timestamp.Format(layout),
			)
			items = append(items, itemString)
		}

	case FilterTodos:
		for id, person := range peopleDB {
			itemString := fmt.Sprintf("%s (%s) - Telefone: %s",
				person.Name,
				id,
				person.Phone,
			)
			items = append(items, itemString)
		}

	case FilterUltimaHora:
		umaHoraAtras := now.Add(-1 * time.Hour)
		for _, entry := range registryLog {
			if entry.TimestampIn.After(umaHoraAtras) ||
				(!entry.TimestampOut.IsZero() && entry.TimestampOut.After(umaHoraAtras)) {

				items = append(items, formatLogEntry(entry, layout))
			}
		}

	case FilterDia:
		ano, mes, dia := now.Date()
		inicioDoDia := time.Date(ano, mes, dia, 0, 0, 0, 0, now.Location())
		for _, entry := range registryLog {
			if entry.TimestampIn.After(inicioDoDia) ||
				(!entry.TimestampOut.IsZero() && entry.TimestampOut.After(inicioDoDia)) {

				items = append(items, formatLogEntry(entry, layout))
			}
		}

	case FilterSaidas:
		for _, entry := range registryLog {
			if !entry.TimestampOut.IsZero() {
				items = append(items, formatLogEntry(entry, layout))
			}
		}

	case FilterCompleto:
		fallthrough
	default:
		for _, entry := range registryLog {
			items = append(items, formatLogEntry(entry, layout))
		}
	}

	list.Set(items)
}

func clearFields(entries ...*widget.Entry) {
	for _, entry := range entries {
		entry.SetText("")
	}
}

func saveData() {
	data := AppData{
		People:   peopleDB,
		Registry: registryLog,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Erro ao converter dados para JSON:", err)
		return
	}

	err = os.WriteFile(dataFile, jsonData, 0644)
	if err != nil {
		fmt.Println("Erro ao salvar arquivo JSON:", err)
	}
}

func loadData() {
	jsonData, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Arquivo de dados não encontrado. Iniciando com estado vazio.")
			return
		}
		fmt.Println("Erro ao ler arquivo JSON:", err)
		return
	}

	var data AppData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Erro ao converter JSON para dados:", err)
		return
	}

	peopleDB = data.People
	registryLog = data.Registry

	for i, entry := range registryLog {
		if entry.TimestampOut.IsZero() {
			activeEntries[entry.PersonID] = i
		}
	}
}
