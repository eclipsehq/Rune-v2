package web

import (
	"encoding/json"
	"fmt"
	"os"
	"net/http"
	"github.com/bwmarrin/discordgo"
	"rune/internal/cmds"
	"rune/internal/config"
	"sort"
	"strings"
)

var botSession *discordgo.Session

var OnTokenChange func(string)

func SetSession(s *discordgo.Session) {
	botSession = s
}

func Start() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/config", configHandler)
	http.HandleFunc("/commands", commandsHandler)
	http.HandleFunc("/presence", presenceHandler)
	http.HandleFunc("/stop", stopHandler)

	go func() {
		fmt.Println("[WEB] Dashboard initialized on http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("[WEB] Server error: %v\n", err)
		}
	}()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	config.Mu.Lock()
	prefix := config.Cfg.Prefix
	ownerID := config.Cfg.OwnerID
	config.Mu.Unlock()

	content := fmt.Sprintf(`
		<div class="stats-grid">
			<div class="stat-card"><div class="stat-label">Prefix</div><div class="stat-val">%s</div></div>
			<div class="stat-card"><div class="stat-label">Owner ID</div><div class="stat-val" style="font-size:16px">%s</div></div>
			<div class="stat-card"><div class="stat-label">Total Commands</div><div class="stat-val">%d</div></div>
		</div>
		<div class="config-box" style="margin-top:20px; animation: fadeIn 0.8s ease-out;">
			<h2 class="section-title">Console Overview</h2>
			<p style="color:var(--text-dim); line-height:1.8; font-size: 15px;">
				Rune V2 is active and running. This dashboard provides a central interface for managing your self-bot's behavior, 
				modifying account presence, and configuring command aliases on the fly.
			</p>
			<div style="display:grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-top: 30px;">
				<a href="/config" class="nav-item" style="background:var(--border); text-align:center; color:var(--accent);">Adjust Settings</a>
				<a href="/commands" class="nav-item" style="background:var(--border); text-align:center; color:var(--accent);">Browse Library</a>
			</div>
		</div>
	`, prefix, ownerID, len(cmds.Commands))

	w.Write([]byte(renderLayout("Overview", content, "index")))
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		oldToken := config.Cfg.Token
		newToken := r.FormValue("token")

		config.Mu.Lock()
		config.Cfg.Token = newToken
		config.Cfg.OwnerID = r.FormValue("owner_id")
		config.Cfg.Prefix = r.FormValue("prefix")
		data, _ := json.MarshalIndent(config.Cfg, "", "  ")
		_ = os.MkdirAll("cfg", 0755)
		_ = os.WriteFile("cfg/config.json", data, 0644)
		config.Mu.Unlock()

		if newToken != oldToken && OnTokenChange != nil {
			OnTokenChange(newToken)
		}
		http.Redirect(w, r, "/config", http.StatusSeeOther)
		return
	}

	config.Mu.Lock()
	token := config.Cfg.Token
	prefix := config.Cfg.Prefix
	ownerID := config.Cfg.OwnerID
	config.Mu.Unlock()

	form := fmt.Sprintf(`
		<h2 class="section-title">Global Configuration</h2>
		<div class="config-box">
			<form method="POST">
				<div style="margin-bottom:25px">
					<label class="input-label">Access Token</label>
					<input type="password" name="token" value="%s" placeholder="Bot/User Token">
				</div>
				<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 30px; margin-bottom: 30px;">
					<div>
						<label class="input-label">Owner ID</label>
						<input type="text" name="owner_id" value="%s">
					</div>
					<div>
						<label class="input-label">Prefix</label>
						<input type="text" name="prefix" value="%s">
					</div>
				</div>
				<button type="submit" class="save-btn">Apply Settings</button>
			</form>
		</div>
	`, token, ownerID, prefix)

	w.Write([]byte(renderLayout("Settings", form, "config")))
}

func commandsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		cmdName := r.FormValue("cmd_name")
		newAliasesRaw := r.FormValue("aliases")
		var cleaned []string
		for _, a := range strings.Split(newAliasesRaw, ",") {
			if trimmed := strings.TrimSpace(a); trimmed != "" {
				cleaned = append(cleaned, trimmed)
			}
		}
		if cmd, ok := cmds.Commands[cmdName]; ok {
			cmd.Aliases = cleaned
			cmds.Commands[cmdName] = cmd
		}
		http.Redirect(w, r, "/commands", http.StatusSeeOther)
		return
	}

	config.Mu.Lock()
	prefix := config.Cfg.Prefix
	config.Mu.Unlock()

	var cmdHtml strings.Builder
	keys := make([]string, 0, len(cmds.Commands))
	for k := range cmds.Commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		cmd := cmds.Commands[name]
		aliases := strings.Join(cmd.Aliases, ", ")
		cmdHtml.WriteString(fmt.Sprintf(`
			<div class="cmd-card" onclick="toggleDetails('%s')">
				<div class="cmd-header">
					<span class="cmd-name">%s%s</span>
					<span class="cmd-cat">%s</span>
				</div>
				<p class="cmd-desc">%s</p>
				<div id="details-%s" class="cmd-details">
					<div class="details-inner">
						<div class="detail-item"><strong>Category:</strong> %s</div>
						<div class="detail-item"><strong>Aliases:</strong> %s</div>
						<div class="detail-divider" style="height:1px; background:var(--border); margin:15px 0;"></div>
						<form method="POST" onclick="event.stopPropagation()">
							<input type="hidden" name="cmd_name" value="%s">
							<label class="input-label" style="font-size: 10px;">Quick Alias Update</label>
							<div style="display:flex; gap:10px; margin-top:8px;">
								<input type="text" name="aliases" value="%s" placeholder="New aliases..." class="alias-input">
								<button type="submit" class="alias-btn">Save</button>
							</div>
						</form>
					</div>
				</div>
			</div>`, name, prefix, name, strings.ToUpper(cmd.Category), cmd.Description, name, cmd.Category, aliases, name, aliases))
	}

	content := fmt.Sprintf(`
		<h2 class="section-title">Command Library</h2>
		<div style="margin-bottom: 30px;">
			<input type="text" id="cmdSearch" onkeyup="filterCommands()" placeholder="Search commands by name or category..." style="width:100%%; padding:16px; background:#080808; border:1px solid var(--border); border-radius:14px; color:#fff; outline:none; font-family:inherit;">
		</div>
		<div class="cmd-grid">%s</div>
		<script>
			function toggleDetails(name) {
				document.getElementById('details-' + name).classList.toggle('active');
			}
			function filterCommands() {
				let input = document.getElementById('cmdSearch').value.toLowerCase();
				let cards = document.getElementsByClassName('cmd-card');
				for (let card of cards) {
					let name = card.querySelector('.cmd-name').innerText.toLowerCase();
					let cat = card.querySelector('.cmd-cat').innerText.toLowerCase();
					if (name.includes(input) || cat.includes(input)) {
						card.style.display = "";
					} else {
						card.style.display = "none";
					}
				}
			}
		</script>
	`, cmdHtml.String())

	w.Write([]byte(renderLayout("Commands", content, "commands")))
}

func presenceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && botSession != nil {
		status := r.FormValue("status")
		if status != "" {
			_ = botSession.UpdateStatusComplex(discordgo.UpdateStatusData{
				Status: status,
			})
		}
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		os.Exit(0)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func getProfileData() (avatar string, statusColor string, username string, status string) {
	status = "Offline"
	statusColor = "#ff4444"
	username = "Not Connected"
	avatar = "https://cdn.discordapp.com/embed/avatars/0.png"

	if botSession != nil && botSession.State != nil && botSession.State.User != nil {
		u := botSession.State.User
		status = "Connected"
		statusColor = "#00ffcc"
		username = fmt.Sprintf("%s#%s", u.Username, u.Discriminator)
		avatar = u.AvatarURL("128")
	}
	return
}

func renderLayout(pageTitle string, body string, active string) string {
	avatar, statusColor, username, statusText := getProfileData()

	navActive := map[string]string{"index": "", "config": "", "commands": ""}
	navActive[active] = "active"

	tpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Rune V2 | %s</title>
    <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;600;800&display=swap" rel="stylesheet">
    <style>
        :root { --accent: #00ffcc; --bg: #050505; --surface: #0a0a0a; --card: #111111; --border: #1a1a1a; --text: #f0f0f0; --text-dim: #777; --shadow: 0 10px 40px rgba(0,0,0,0.6); }
        * { transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1); box-sizing: border-box; }
        body { background: var(--bg); color: var(--text); font-family: 'Plus Jakarta Sans', sans-serif; margin: 0; display: flex; min-height: 100vh; }
        .sidebar { width: 320px; background: var(--surface); border-right: 1px solid var(--border); padding: 50px 30px; display: flex; flex-direction: column; position: sticky; top: 0; height: 100vh; animation: slideIn 0.6s ease-out; }
        .main-content { flex: 1; padding: 60px; height: 100vh; overflow-y: auto; scroll-behavior: smooth; }
        .logo { font-weight: 800; font-size: 28px; color: var(--accent); letter-spacing: -1.5px; margin-bottom: 50px; }
        .profile-card { background: #141414; border: 1px solid var(--border); padding: 25px; border-radius: 24px; margin-bottom: 30px; box-shadow: var(--shadow); }
        .pfp-wrap { position: relative; width: 64px; height: 64px; margin-bottom: 15px; }
        .pfp { width: 100%%; height: 100%%; border-radius: 18px; border: 2px solid %s; }
        .status-dot { position: absolute; bottom: -2px; right: -2px; width: 14px; height: 14px; border-radius: 50%%; background: %s; border: 3px solid var(--card); }
        .user-name { font-weight: 700; font-size: 18px; margin-bottom: 4px; }
        .user-status { font-size: 12px; color: var(--text-dim); margin-bottom: 15px; }
        .presence-select { width: 100%%; background: #080808; color: var(--text); border: 1px solid var(--border); padding: 10px; border-radius: 10px; font-size: 12px; outline: none; cursor: pointer; }
        .nav-links { display: flex; flex-direction: column; gap: 8px; flex-grow: 1; }
        .nav-item { padding: 14px 20px; border-radius: 14px; text-decoration: none; color: var(--text-dim); font-weight: 600; font-size: 14px; }
        .nav-item:hover, .nav-item.active { background: rgba(0,255,204,0.05); color: var(--accent); }
        .repo-btn { background: var(--accent); color: #000; padding: 16px; border-radius: 14px; text-decoration: none; font-weight: 800; font-size: 14px; text-align: center; margin-top: 20px; }
        .stats-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 25px; margin-bottom: 45px; }
        .stat-card { background: var(--card); border: 1px solid var(--border); padding: 30px; border-radius: 20px; box-shadow: var(--shadow); animation: fadeIn 0.6s ease-out; }
        .stat-label { color: var(--text-dim); font-size: 11px; text-transform: uppercase; font-weight: 800; margin-bottom: 10px; letter-spacing: 1px; }
        .stat-val { font-size: 26px; font-weight: 800; color: #fff; }
        .section-title { font-size: 26px; font-weight: 800; margin-bottom: 35px; display: flex; align-items: center; gap: 15px; animation: fadeIn 0.8s ease-out; }
        .section-title::after { content: ''; flex: 1; height: 1px; background: var(--border); }
        .config-box { background: var(--card); border: 1px solid var(--border); padding: 45px; border-radius: 28px; margin-bottom: 60px; box-shadow: var(--shadow); animation: fadeIn 1s ease-out; }
        .input-label { display: block; font-size: 11px; color: var(--text-dim); font-weight: 800; text-transform: uppercase; margin-bottom: 10px; letter-spacing: 0.5px; }
        input { width: 100%%; background: #080808; border: 1px solid var(--border); padding: 16px; border-radius: 14px; color: #fff; font-family: inherit; outline: none; }
        input:focus { border-color: var(--accent); box-shadow: 0 0 20px rgba(0,255,204,0.1); }
        .save-btn { background: var(--accent); color: #000; border: none; padding: 18px; border-radius: 14px; font-weight: 800; cursor: pointer; width: 100%%; margin-top: 10px; }
        .stop-btn { background: #ff444415; color: #ff4444; border: 1px solid #ff444433; padding: 16px; border-radius: 14px; font-weight: 700; cursor: pointer; width: 100%%; }
        .stop-btn:hover { background: #ff4444; color: #fff; }
        .cmd-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(360px, 1fr)); gap: 25px; animation: fadeIn 1.2s ease-out; }
        .cmd-card { background: var(--card); border: 1px solid var(--border); padding: 32px; border-radius: 24px; cursor: pointer; box-shadow: 0 4px 20px rgba(0,0,0,0.2); }
        .cmd-card:hover { border-color: var(--accent); transform: translateY(-6px); box-shadow: 0 12px 30px rgba(0,0,0,0.4); }
        .cmd-name { font-weight: 800; color: var(--accent); font-size: 20px; }
        .cmd-cat { font-size: 10px; background: rgba(255,255,255,0.05); padding: 5px 12px; border-radius: 8px; color: #aaa; font-weight: 800; }
        .cmd-desc { color: var(--text-dim); font-size: 14px; line-height: 1.6; margin: 15px 0 0 0; }
        .cmd-details { max-height: 0; overflow: hidden; opacity: 0; transition: all 0.4s ease; }
        .cmd-details.active { max-height: 400px; opacity: 1; margin-top: 25px; padding-top: 25px; border-top: 1px solid var(--border); }
        .alias-input { background: #080808; border: 1px solid var(--border); color: #fff; font-size: 12px; padding: 12px; border-radius: 10px; flex: 1; }
        .alias-btn { background: #222; color: #fff; border: none; padding: 0 20px; border-radius: 10px; font-weight: 700; cursor: pointer; }
        @keyframes slideIn { from { transform: translateX(-40px); opacity: 0; } to { transform: translateX(0); opacity: 1; } }
        @keyframes fadeIn { from { transform: translateY(30px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
        ::-webkit-scrollbar { width: 6px; }
        ::-webkit-scrollbar-track { background: var(--bg); }
        ::-webkit-scrollbar-thumb { background: var(--border); border-radius: 10px; }
    </style>
</head>
<body>
    <div class="sidebar">
        <div class="logo">RUNE / V2</div>        
        <div class="profile-card">
            <div class="pfp-wrap"><img src="%s" class="pfp"><div class="status-dot"></div></div>
            <div class="user-name">%s</div><div class="user-status">%s</div>
            <form method="POST" action="/presence" style="margin-top:15px">
                <select name="status" onchange="this.form.submit()" class="presence-select">
                    <option value="" disabled selected>Update Presence</option>
                    <option value="online">Online</option>
                    <option value="idle">Idle</option>
                    <option value="dnd">DND</option>
                    <option value="invisible">Invisible</option>
                </select>
            </form>
        </div>
        <div class="nav-links">
            <a href="/" class="nav-item %s">Overview</a>
            <a href="/config" class="nav-item %s">Configuration</a>
            <a href="/commands" class="nav-item %s">Command Library</a>
        </div>
        <form method="POST" action="/stop" style="margin-top: auto;">
            <button type="submit" class="stop-btn">Shutdown Bot</button>
        </form>
    </div>
    <div class="main-content">
        [[BODY]]
    </div>
</body>
</html>`

	html := fmt.Sprintf(tpl, pageTitle, statusColor, statusColor, avatar, username, statusText, navActive["index"], navActive["config"], navActive["commands"])
	return strings.Replace(html, "[[BODY]]", body, 1)
}
