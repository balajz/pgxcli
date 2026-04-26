package config

type SyntaxHighlightStyle string

const (
	SyntaxStyleDefault             SyntaxHighlightStyle = "default"
	SyntaxStyleABAP                SyntaxHighlightStyle = "abap"
	SyntaxStyleAlgol               SyntaxHighlightStyle = "algol"
	SyntaxStyleAlgolNU             SyntaxHighlightStyle = "algol_nu"
	SyntaxStyleArduino             SyntaxHighlightStyle = "arduino"
	SyntaxStyleAshen               SyntaxHighlightStyle = "ashen"
	SyntaxStyleAuraThemeDark       SyntaxHighlightStyle = "aura-theme-dark"
	SyntaxStyleAuraThemeDarkSoft   SyntaxHighlightStyle = "aura-theme-dark-soft"
	SyntaxStyleAutumn              SyntaxHighlightStyle = "autumn"
	SyntaxStyleAverage             SyntaxHighlightStyle = "average"
	SyntaxStyleBase16Snazzy        SyntaxHighlightStyle = "base16-snazzy"
	SyntaxStyleBorland             SyntaxHighlightStyle = "borland"
	SyntaxStyleBW                  SyntaxHighlightStyle = "bw"
	SyntaxStyleCatppuccinFrappe    SyntaxHighlightStyle = "catppuccin-frappe"
	SyntaxStyleCatppuccinLatte     SyntaxHighlightStyle = "catppuccin-latte"
	SyntaxStyleCatppuccinMacchiato SyntaxHighlightStyle = "catppuccin-macchiato"
	SyntaxStyleCatppuccinMocha     SyntaxHighlightStyle = "catppuccin-mocha"
	SyntaxStyleColorful            SyntaxHighlightStyle = "colorful"
	SyntaxStyleDarcula             SyntaxHighlightStyle = "darcula" //nolint:misspell
	SyntaxStyleDoomOne             SyntaxHighlightStyle = "doom-one"
	SyntaxStyleDoomOne2            SyntaxHighlightStyle = "doom-one2"
	SyntaxStyleDracula             SyntaxHighlightStyle = "dracula"
	SyntaxStyleEmacs               SyntaxHighlightStyle = "emacs"
	SyntaxStyleEvergarden          SyntaxHighlightStyle = "evergarden"
	SyntaxStyleFriendly            SyntaxHighlightStyle = "friendly"
	SyntaxStyleFruity              SyntaxHighlightStyle = "fruity"
	SyntaxStyleGitHub              SyntaxHighlightStyle = "github"
	SyntaxStyleGitHubDark          SyntaxHighlightStyle = "github-dark"
	SyntaxStyleGruvbox             SyntaxHighlightStyle = "gruvbox"
	SyntaxStyleGruvboxLight        SyntaxHighlightStyle = "gruvbox-light"
	SyntaxStyleHRHighContrast      SyntaxHighlightStyle = "hr_high_contrast"
	SyntaxStyleHRDark              SyntaxHighlightStyle = "hrdark"
	SyntaxStyleIgor                SyntaxHighlightStyle = "igor"
	SyntaxStyleKanagawaDragon      SyntaxHighlightStyle = "kanagawa-dragon"
	SyntaxStyleKanagawaLotus       SyntaxHighlightStyle = "kanagawa-lotus"
	SyntaxStyleKanagawaWave        SyntaxHighlightStyle = "kanagawa-wave"
	SyntaxStyleLovelace            SyntaxHighlightStyle = "lovelace"
	SyntaxStyleManni               SyntaxHighlightStyle = "manni"
	SyntaxStyleModusOperandi       SyntaxHighlightStyle = "modus-operandi"
	SyntaxStyleModusVivendi        SyntaxHighlightStyle = "modus-vivendi"
	SyntaxStyleMonokai             SyntaxHighlightStyle = "monokai"
	SyntaxStyleMonokaiLight        SyntaxHighlightStyle = "monokailight"
	SyntaxStyleMurphy              SyntaxHighlightStyle = "murphy"
	SyntaxStyleNative              SyntaxHighlightStyle = "native"
	SyntaxStyleNord                SyntaxHighlightStyle = "nord"
	SyntaxStyleNordic              SyntaxHighlightStyle = "nordic"
	SyntaxStyleOnedark             SyntaxHighlightStyle = "onedark"
	SyntaxStyleOneSEnterprise      SyntaxHighlightStyle = "onesenterprise"
	SyntaxStyleParaisoDark         SyntaxHighlightStyle = "paraiso-dark"
	SyntaxStyleParaisoLight        SyntaxHighlightStyle = "paraiso-light"
	SyntaxStylePastie              SyntaxHighlightStyle = "pastie"
	SyntaxStylePerldoc             SyntaxHighlightStyle = "perldoc"
	SyntaxStylePygments            SyntaxHighlightStyle = "pygments"
	SyntaxStyleRainbowDash         SyntaxHighlightStyle = "rainbow_dash"
	SyntaxStyleRosePine            SyntaxHighlightStyle = "rose-pine"
	SyntaxStyleRosePineDawn        SyntaxHighlightStyle = "rose-pine-dawn"
	SyntaxStyleRosePineMoon        SyntaxHighlightStyle = "rose-pine-moon"
	SyntaxStyleRPGLE               SyntaxHighlightStyle = "RPGLE"
	SyntaxStyleRrt                 SyntaxHighlightStyle = "rrt"
	SyntaxStyleSolarizedDark       SyntaxHighlightStyle = "solarized-dark"
	SyntaxStyleSolarizedDark256    SyntaxHighlightStyle = "solarized-dark256"
	SyntaxStyleSolarizedLight      SyntaxHighlightStyle = "solarized-light"
	SyntaxStyleSwapoff             SyntaxHighlightStyle = "swapoff"
	SyntaxStyleTango               SyntaxHighlightStyle = "tango"
	SyntaxStyleTokyoNightDay       SyntaxHighlightStyle = "tokyonight-day"
	SyntaxStyleTokyoNightMoon      SyntaxHighlightStyle = "tokyonight-moon"
	SyntaxStyleTokyoNightNight     SyntaxHighlightStyle = "tokyonight-night"
	SyntaxStyleTokyoNightStorm     SyntaxHighlightStyle = "tokyonight-storm"
	SyntaxStyleTrac                SyntaxHighlightStyle = "trac"
	SyntaxStyleVim                 SyntaxHighlightStyle = "vim"
	SyntaxStyleVS                  SyntaxHighlightStyle = "vs"
	SyntaxStyleVulcan              SyntaxHighlightStyle = "vulcan"
	SyntaxStyleWitchhazel          SyntaxHighlightStyle = "witchhazel"
	SyntaxStyleXcode               SyntaxHighlightStyle = "xcode"
	SyntaxStyleXcodeDark           SyntaxHighlightStyle = "xcode-dark"
)

func (s SyntaxHighlightStyle) isValid() bool {
	_, ok := validSyntaxHighlightStyles[s]
	return ok
}

var validSyntaxHighlightStyles = map[SyntaxHighlightStyle]struct{}{
	SyntaxStyleDefault:             {},
	SyntaxStyleABAP:                {},
	SyntaxStyleAlgol:               {},
	SyntaxStyleAlgolNU:             {},
	SyntaxStyleArduino:             {},
	SyntaxStyleAshen:               {},
	SyntaxStyleAuraThemeDark:       {},
	SyntaxStyleAuraThemeDarkSoft:   {},
	SyntaxStyleAutumn:              {},
	SyntaxStyleAverage:             {},
	SyntaxStyleBase16Snazzy:        {},
	SyntaxStyleBorland:             {},
	SyntaxStyleBW:                  {},
	SyntaxStyleCatppuccinFrappe:    {},
	SyntaxStyleCatppuccinLatte:     {},
	SyntaxStyleCatppuccinMacchiato: {},
	SyntaxStyleCatppuccinMocha:     {},
	SyntaxStyleColorful:            {},
	SyntaxStyleDarcula:             {},
	SyntaxStyleDoomOne:             {},
	SyntaxStyleDoomOne2:            {},
	SyntaxStyleDracula:             {},
	SyntaxStyleEmacs:               {},
	SyntaxStyleEvergarden:          {},
	SyntaxStyleFriendly:            {},
	SyntaxStyleFruity:              {},
	SyntaxStyleGitHub:              {},
	SyntaxStyleGitHubDark:          {},
	SyntaxStyleGruvbox:             {},
	SyntaxStyleGruvboxLight:        {},
	SyntaxStyleHRHighContrast:      {},
	SyntaxStyleHRDark:              {},
	SyntaxStyleIgor:                {},
	SyntaxStyleKanagawaDragon:      {},
	SyntaxStyleKanagawaLotus:       {},
	SyntaxStyleKanagawaWave:        {},
	SyntaxStyleLovelace:            {},
	SyntaxStyleManni:               {},
	SyntaxStyleModusOperandi:       {},
	SyntaxStyleModusVivendi:        {},
	SyntaxStyleMonokai:             {},
	SyntaxStyleMonokaiLight:        {},
	SyntaxStyleMurphy:              {},
	SyntaxStyleNative:              {},
	SyntaxStyleNord:                {},
	SyntaxStyleNordic:              {},
	SyntaxStyleOnedark:             {},
	SyntaxStyleOneSEnterprise:      {},
	SyntaxStyleParaisoDark:         {},
	SyntaxStyleParaisoLight:        {},
	SyntaxStylePastie:              {},
	SyntaxStylePerldoc:             {},
	SyntaxStylePygments:            {},
	SyntaxStyleRainbowDash:         {},
	SyntaxStyleRosePine:            {},
	SyntaxStyleRosePineDawn:        {},
	SyntaxStyleRosePineMoon:        {},
	SyntaxStyleRPGLE:               {},
	SyntaxStyleRrt:                 {},
	SyntaxStyleSolarizedDark:       {},
	SyntaxStyleSolarizedDark256:    {},
	SyntaxStyleSolarizedLight:      {},
	SyntaxStyleSwapoff:             {},
	SyntaxStyleTango:               {},
	SyntaxStyleTokyoNightDay:       {},
	SyntaxStyleTokyoNightMoon:      {},
	SyntaxStyleTokyoNightNight:     {},
	SyntaxStyleTokyoNightStorm:     {},
	SyntaxStyleTrac:                {},
	SyntaxStyleVim:                 {},
	SyntaxStyleVS:                  {},
	SyntaxStyleVulcan:              {},
	SyntaxStyleWitchhazel:          {},
	SyntaxStyleXcode:               {},
	SyntaxStyleXcodeDark:           {},
}

// OnErrorAction controls multi-statement behavior after a statement fails.
type OnErrorAction string

const (
	// OnErrorResume continues executing remaining statements after an error.
	OnErrorResume OnErrorAction = "RESUME"
	// OnErrorStop stops execution immediately when a statement fails.
	OnErrorStop OnErrorAction = "STOP"
)

func (a OnErrorAction) isValid() bool {
	switch a {
	case OnErrorResume, OnErrorStop:
		return true
	default:
		return false
	}
}

type TableColor string

const (
	FgDefault   TableColor = "default" // default is cyan for header, white for column and caption
	FgBlack     TableColor = "black"
	FgRed       TableColor = "red"
	FgGreen     TableColor = "green"
	FgYellow    TableColor = "yellow"
	FgBlue      TableColor = "blue"
	FgMagenta   TableColor = "magenta"
	FgCyan      TableColor = "cyan"
	FgWhite     TableColor = "white"
	FgHiBlack   TableColor = "black+"
	FgHiRed     TableColor = "red+"
	FgHiGreen   TableColor = "green+"
	FgHiYellow  TableColor = "yellow+"
	FgHiBlue    TableColor = "blue+"
	FgHiMagenta TableColor = "magenta+"
	FgHiCyan    TableColor = "cyan+"
	FgHiWhite   TableColor = "white+"
)

func (c TableColor) isValid() bool {
	_, ok := validTableColors[c]
	return ok
}

var validTableColors = map[TableColor]struct{}{
	FgDefault:   {},
	FgBlack:     {},
	FgRed:       {},
	FgGreen:     {},
	FgYellow:    {},
	FgBlue:      {},
	FgMagenta:   {},
	FgCyan:      {},
	FgWhite:     {},
	FgHiBlack:   {},
	FgHiRed:     {},
	FgHiGreen:   {},
	FgHiYellow:  {},
	FgHiBlue:    {},
	FgHiMagenta: {},
	FgHiCyan:    {},
	FgHiWhite:   {},
}

type TableStyle string

const (
	StyleNone        TableStyle = "none"
	StyleASCII       TableStyle = "ascii"
	StyleLight       TableStyle = "light"
	StyleHeavy       TableStyle = "heavy"
	StyleDouble      TableStyle = "double"
	StyleDoubleLong  TableStyle = "double_long"
	StyleLightHeavy  TableStyle = "light_heavy"
	StyleHeavyLight  TableStyle = "heavy_light"
	StyleLightDouble TableStyle = "light_double"
	StyleDoubleLight TableStyle = "double_light"
	StyleRounded     TableStyle = "rounded"
	StyleMarkdown    TableStyle = "markdown"
	StyleGraphical   TableStyle = "graphical"
	StyleMerger      TableStyle = "merger"
	StyleDefault     TableStyle = "default"
	StyleDotted      TableStyle = "dotted"
	StyleArrow       TableStyle = "arrow"
	StyleStarry      TableStyle = "starry"
	StyleHearts      TableStyle = "hearts"
	StyleCircuit     TableStyle = "circuit"
	StyleNature      TableStyle = "nature"
	StyleArtistic    TableStyle = "artistic"
	Style8Bit        TableStyle = "8bit"
	StyleChaos       TableStyle = "chaos"
	StyleDots        TableStyle = "dots"
	StyleBlocks      TableStyle = "blocks"
	StyleZen         TableStyle = "zen"
	StyleVintage     TableStyle = "vintage"
	StyleSketch      TableStyle = "sketch"
	StyleArrowDouble TableStyle = "arrow_double"
	StyleCelestial   TableStyle = "celestial"
	StyleCyber       TableStyle = "cyber"
	StyleRunic       TableStyle = "runic"
	StyleIndustrial  TableStyle = "industrial"
	StyleInk         TableStyle = "ink"
	StyleArcade      TableStyle = "arcade"
	StyleBlossom     TableStyle = "blossom"
	StyleFrosted     TableStyle = "frosted"
	StyleMosaic      TableStyle = "mosaic"
	StyleUFO         TableStyle = "ufo"
	StyleSteampunk   TableStyle = "steampunk"
	StyleGalaxy      TableStyle = "galaxy"
	StyleJazz        TableStyle = "jazz"
	StylePuzzle      TableStyle = "puzzle"
	StyleHypno       TableStyle = "hypno"
)

func (s TableStyle) isValid() bool {
	_, ok := validTableStyles[s]
	return ok
}

var validTableStyles = map[TableStyle]struct{}{
	StyleNone:        {},
	StyleASCII:       {},
	StyleLight:       {},
	StyleHeavy:       {},
	StyleDouble:      {},
	StyleDoubleLong:  {},
	StyleLightHeavy:  {},
	StyleHeavyLight:  {},
	StyleLightDouble: {},
	StyleDoubleLight: {},
	StyleRounded:     {},
	StyleMarkdown:    {},
	StyleGraphical:   {},
	StyleMerger:      {},
	StyleDefault:     {},
	StyleDotted:      {},
	StyleArrow:       {},
	StyleStarry:      {},
	StyleHearts:      {},
	StyleCircuit:     {},
	StyleNature:      {},
	StyleArtistic:    {},
	Style8Bit:        {},
	StyleChaos:       {},
	StyleDots:        {},
	StyleBlocks:      {},
	StyleZen:         {},
	StyleVintage:     {},
	StyleSketch:      {},
	StyleArrowDouble: {},
	StyleCelestial:   {},
	StyleCyber:       {},
	StyleRunic:       {},
	StyleIndustrial:  {},
	StyleInk:         {},
	StyleArcade:      {},
	StyleBlossom:     {},
	StyleFrosted:     {},
	StyleMosaic:      {},
	StyleUFO:         {},
	StyleSteampunk:   {},
	StyleGalaxy:      {},
	StyleJazz:        {},
	StylePuzzle:      {},
	StyleHypno:       {},
}
