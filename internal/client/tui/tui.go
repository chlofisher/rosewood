package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/list"
	"charm.land/lipgloss/v2"

	"chlofisher.com/rosewood/internal/client"
	"chlofisher.com/rosewood/internal/metadata"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	clientInstance *client.Client	

	songs []*metadata.Song

	searchBar textinput.Model
	songList list.Model
}

type songItem struct {
	title, desc string
}

func (i songItem) Title() string       { return i.title }
func (i songItem) Description() string { return i.desc }
func (i songItem) FilterValue() string { return i.title }

func songsToItems(songs []*metadata.Song) []list.Item {
	items := make([]list.Item, len(songs))
	for i, s := range songs {
		items[i] = songItem{title: s.Title, desc: s.Artist}
	}
	return items
}

func New(c *client.Client) tea.Model {
	search := textinput.New()
	search.SetVirtualCursor(true)
	search.SetWidth(30)
	search.Placeholder = "Press / to search..."


	songs := make([]*metadata.Song, 0)
	l := list.New(songsToItems(songs), list.NewDefaultDelegate(), 0, 0)

	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowHelp(true)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()

	return model{
		searchBar: search,
		clientInstance: c,
		songs: songs,
		songList: l,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd	

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "/":
			m.searchBar.Focus()
			return m, textinput.Blink
		case "esc":
			m.searchBar.Blur()
		case "enter":
			idx := m.songList.Cursor()
			m.clientInstance.PlaySongID(m.songs[idx].ID)
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.songList.SetSize(msg.Width-h, msg.Height-v-1)
	}

	m.searchBar, cmd = m.searchBar.Update(msg)

	songs, _ := m.clientInstance.SearchSongs(m.searchBar.Value())
	m.songs = songs
	m.songList.SetItems(songsToItems(songs))

	m.songList, cmd = m.songList.Update(msg)

	return m, cmd
}

func (m model) View() tea.View {
	str := lipgloss.JoinVertical(
		lipgloss.Left,
		m.searchBar.View(),
		m.listView(),
	)
	v := tea.NewView(str)

	var c *tea.Cursor
	if !m.searchBar.VirtualCursor() {
		c = m.searchBar.Cursor()
	}
	v.Cursor = c

	return v
}

func (m model) searchView() string {
	return m.searchBar.View()
}

func (m model) listView() string {
	return docStyle.Render(m.songList.View())
}

