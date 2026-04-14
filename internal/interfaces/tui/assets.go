package tui

// Scaling: Simplified these to fit terminal better (approx 50%)

const KeroChan = `    \ \  _  / /
     \ \/ \/ /   __  __
  _  ( @   @ )  /  \/  \
 / \/ (  ^  )  (  SHINE )
 \_  / \---/    \__  __/
   \/   | |        \/`

const Bersert = `      \__/
     (o  o)
      \  /
      |  |
      |  |`

const Evangelion = `     ||
     ||
    /  \
   |    |
    \__/`

// Actual User Templates (truncated/scaled for TUI use)
var AsciiIcons = []string{
	// Kero-chan (Scaled down)
	`  ⣀⡴⠞⠋⠋⠛⠳⢦⡀⠀⠀⢀⣀⣤⣤⣀⣀⠀⠀⠀⣠⡶⠟⠉⠉⠛⠶⣤
⡼⢋⣴⡶⠶⣶⣤⠀⠀⠙⠶⠛⠉⠁⠀⠀⠀⠀⠈⠙⠳⠞⠁⠀⣠⣶⡿⠿⢷⣌⠳
⡇⡿⠄⠀⠀⠀⣀⡾⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⢿⡀⠀⠀⠀⠀⢺⡦⡟
⢿⡸⣦⠀⠀⢰⠇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⡆⠀⠀⣰⢃⡿
⠀⠉⠳⠮⠷⣜⢷⣄⢸⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡇⣠⡾
⠀⠀⠀⠀⠀⠙⢿⡀⠀⠈⠀⠁⠀⠀⠀⠀⠐⠖⠀⠀⠀⠀⠁⠈⠀⠀⠀⣼⠋⠁
⠀⠀⠀⠀⠀⠀⠙⣷⣍⣉⣉⣋⣉⣋⣉⣁⣴⠟⠁`,

	// Berserk (Scaled)
	`    ⣰⣾⠁  ⢦⣾⣤⠆  ⠻⣧
  ⢠⣼⠏      ⣿⡇      ⠈⢷⣄
 ⢰⣾⣿⡁      ⣿⡇      ⢀⣿⣿⠖
   ⠈⠻⣿⣦⣄    ⣿⡇    ⢀⣴⣿⠟⠁
       ⠙⢿⣿⣦⣿⣧⣾⣿⠟⠁
       ⢀⣴⣿⣿⣿⣿⣿⣿⣦⡀
   ⣠⣴⣿⡿⠋      ⢼⣿      ⠈⢻⣿⣷⣄`,

	// Evangelion (Scaled)
	`     ⢸     ⢸
     ⢸     ⢸
     ⢸     ⣷
     ⢸     ⢸
  ⠤⢧⢹    ⢀⢘⡄⡄
     ⡀⣾    ⢸⢹⠋
  ⢀⢅    ⠱⣐⢔⡪⠃
  ⢸⢸    ⠈⢸⠾⡎⠁
  ⣼⢸⠂    ⣠⠯⡄
 ⠙⠘⠩⡇ ⢼⠟⡏⡏⠁`,
}
