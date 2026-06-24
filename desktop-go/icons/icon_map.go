package icons

// ExtensionIconMap maps file extensions to their icon SVG filename (without path prefix).
// The icon files are located in the embedded "icons/" directory.
var ExtensionIconMap = map[string]string{
	// Go
	".go":    "file_type_go.svg",
	".gomod": "file_type_go_package.svg",
	".gosum": "file_type_go_package.svg",
	".work":  "file_type_go_work.svg",

	// JavaScript
	".js":  "file_type_js.svg",
	".mjs": "file_type_js.svg",
	".cjs": "file_type_js.svg",
	".jsx": "file_type_reactjs.svg",

	// TypeScript
	".ts":   "file_type_typescript.svg",
	".tsx":  "file_type_reactts.svg",
	".d.ts": "file_type_typescriptdef.svg",

	// Python
	".py":    "file_type_python.svg",
	".pyw":   "file_type_python.svg",
	".pyx":   "file_type_cython.svg",
	".ipynb": "file_type_jupyter.svg",

	// Rust
	".rs": "file_type_rust.svg",

	// C / C++
	".c":   "file_type_c.svg",
	".h":   "file_type_cheader.svg",
	".cpp": "file_type_cpp.svg",
	".cc":  "file_type_cpp.svg",
	".cxx": "file_type_cpp.svg",
	".hpp": "file_type_cppheader.svg",
	".hh":  "file_type_cppheader.svg",

	// C#
	".cs":     "file_type_csharp.svg",
	".csproj": "file_type_csproj.svg",

	// Java
	".java":   "file_type_java.svg",
	".jar":    "file_type_jar.svg",
	".class":  "file_type_class.svg",
	".kt":     "file_type_kotlin.svg",
	".kts":    "file_type_kotlin.svg",
	".groovy": "file_type_groovy.svg",
	".gradle": "file_type_gradle.svg",
	".scala":  "file_type_scala.svg",
	".sbt":    "file_type_sbt.svg",

	// Ruby
	".rb":      "file_type_ruby.svg",
	".erb":     "file_type_erb.svg",
	".rbs":     "file_type_ruby.svg",
	".gemspec": "file_type_bundler.svg",
	"Gemfile":  "file_type_bundler.svg",
	"Rakefile": "file_type_rake.svg",

	// PHP
	".php":       "file_type_php.svg",
	".phtml":     "file_type_php.svg",
	".blade.php": "file_type_blade.svg",

	// HTML / Web
	".html":  "file_type_html.svg",
	".htm":   "file_type_html.svg",
	".xhtml": "file_type_html.svg",

	// CSS
	".css":  "file_type_css.svg",
	".scss": "file_type_scss.svg",
	".sass": "file_type_sass.svg",
	".less": "file_type_less.svg",
	".styl": "file_type_stylus.svg",

	// Vue / Svelte / Astro / Angular
	".vue":    "file_type_vue.svg",
	".svelte": "file_type_svelte.svg",
	".astro":  "file_type_astro.svg",

	// React (handled via .jsx/.tsx above)

	// Dart / Flutter
	".dart": "file_type_dartlang.svg",

	// Swift
	".swift": "file_type_swift.svg",

	// Kotlin (already above)

	// Objective-C
	".m":  "file_type_objectivec.svg",
	".mm": "file_type_objectivecpp.svg",

	// Shell / Scripts
	".sh":   "file_type_shell.svg",
	".bash": "file_type_shell.svg",
	".zsh":  "file_type_shell.svg",
	".fish": "file_type_shell.svg",
	".bat":  "file_type_bat.svg",
	".cmd":  "file_type_bat.svg",
	".ps1":  "file_type_powershell.svg",
	".psm1": "file_type_powershell_psm.svg",
	".psd1": "file_type_powershell_psd.svg",

	// Data / Config
	".json":  "file_type_json.svg",
	".json5": "file_type_json5.svg",
	".jsonc": "file_type_json.svg",
	".yaml":  "file_type_yaml.svg",
	".yml":   "file_type_yaml.svg",
	".toml":  "file_type_toml.svg",
	".ini":   "file_type_ini.svg",
	".env":   "file_type_dotenv.svg",
	".xml":   "file_type_xml.svg",
	".xsl":   "file_type_xsl.svg",
	".xsd":   "file_type_xml.svg",
	".csv":   "file_type_excel.svg",

	// Markdown / Docs
	".md":  "file_type_markdown.svg",
	".mdx": "file_type_mdx.svg",
	".rst": "file_type_rest.svg",
	".txt": "file_type_text.svg",
	".tex": "file_type_tex.svg",
	".pdf": "file_type_pdf.svg",

	// SQL / Databases
	".sql":    "file_type_sql.svg",
	".pgsql":  "file_type_pgsql.svg",
	".mysql":  "file_type_mysql.svg",
	".db":     "file_type_db.svg",
	".sqlite": "file_type_sqlite.svg",

	// Docker
	"Dockerfile":         "file_type_docker.svg",
	".dockerfile":        "file_type_docker.svg",
	".dockerignore":      "file_type_docker.svg",
	"docker-compose.yml": "file_type_docker.svg",

	// GraphQL
	".graphql": "file_type_graphql.svg",
	".gql":     "file_type_graphql.svg",

	// Lua
	".lua": "file_type_lua.svg",

	// Elixir / Erlang
	".ex":   "file_type_elixir.svg",
	".exs":  "file_type_elixir.svg",
	".eex":  "file_type_eex.svg",
	".heex": "file_type_eex.svg",
	".erl":  "file_type_erlang.svg",
	".hrl":  "file_type_erlang.svg",

	// Haskell / PureScript
	".hs":   "file_type_haskell.svg",
	".lhs":  "file_type_haskell.svg",
	".purs": "file_type_purescript.svg",

	// OCaml
	".ml":  "file_type_ocaml.svg",
	".mli": "file_type_ocaml_intf.svg",

	// F#
	".fs":     "file_type_fsharp.svg",
	".fsi":    "file_type_fsharp.svg",
	".fsx":    "file_type_fsharp.svg",
	".fsproj": "file_type_fsproj.svg",

	// Ruby on Rails
	".rhtml": "file_type_rails.svg",

	// Terraform / HCL
	".tf":     "file_type_terraform.svg",
	".tfvars": "file_type_terraform.svg",
	".hcl":    "file_type_hashicorp.svg",

	// Protobuf
	".proto": "file_type_protobuf.svg",

	// WASM
	".wasm": "file_type_wasm.svg",
	".wat":  "file_type_wasm.svg",

	// Assembly
	".asm":  "file_type_assembly.svg",
	".s":    "file_type_assembly.svg",
	".nasm": "file_type_assembly.svg",

	// R
	".r":     "file_type_r.svg",
	".R":     "file_type_r.svg",
	".rmd":   "file_type_rmd.svg",
	".rproj": "file_type_rproj.svg",

	// Julia
	".jl": "file_type_julia.svg",

	// Nim
	".nim":    "file_type_nim.svg",
	".nimble": "file_type_nimble.svg",

	// Zig
	".zig": "file_type_zig.svg",

	// V (Vlang)
	".v": "file_type_vlang.svg",

	// Clojure / ClojureScript
	".clj":  "file_type_clojure.svg",
	".cljs": "file_type_clojurescript.svg",
	".cljc": "file_type_clojure.svg",
	".edn":  "file_type_clojure.svg",

	// Dart (handled above)

	// Solidity
	".sol": "file_type_solidity.svg",

	// Crystal
	".cr": "file_type_crystal.svg",

	// Groovy (handled above)

	// Perl
	".pl":  "file_type_perl.svg",
	".pm":  "file_type_perl.svg",
	".pod": "file_type_perl.svg",
	".p6":  "file_type_perl6.svg",

	// MATLAB
	".matla": "file_type_matlab.svg", // conflict with .m Obj-C — see GetIcon logic
	".mat":   "file_type_matlab.svg",

	// Fortran
	".f":   "file_type_fortran.svg",
	".f90": "file_type_fortran.svg",
	".f95": "file_type_fortran.svg",
	".for": "file_type_fortran.svg",

	// COBOL
	".cob": "file_type_cobol.svg",
	".cbl": "file_type_cobol.svg",

	// Ada
	".ada": "file_type_ada.svg",
	".adb": "file_type_ada.svg",
	".ads": "file_type_ada.svg",

	// Lisp / Racket / Scheme
	".lisp": "file_type_lisp.svg",
	".lsp":  "file_type_lisp.svg",
	".rkt":  "file_type_racket.svg",
	".scm":  "file_type_lisp.svg",

	// Prolog
	".prolog": "file_type_prolog.svg", // duplicated with Perl — see GetIcon logic
	".pro":    "file_type_prolog.svg",

	// Tcl
	".tcl": "file_type_tcl.svg",

	// CoffeeScript
	".coffee": "file_type_coffeescript.svg",

	// D Language
	".d": "file_type_dlang.svg",

	// Vala
	".vala": "file_type_vala.svg",
	".vapi": "file_type_vapi.svg",

	// Mojo
	".mojo": "file_type_mojo.svg",
	".🔥":    "file_type_mojo.svg",

	// Gleam
	".gleam": "file_type_gleam.svg",

	// Nix
	".nix": "file_type_nix.svg",

	// Dhall
	".dhall": "file_type_dhall.svg",

	// Idris
	".idr": "file_type_idris.svg",

	// Agda
	".agda": "file_type_agda.svg",

	// Lean
	".lean": "file_type_lean.svg",

	// GDScript (Godot)
	".gd": "file_type_gdscript.svg",

	// HLSL / GLSL / Metal
	".hlsl":  "file_type_hlsl.svg",
	".glsl":  "file_type_glsl.svg",
	".vert":  "file_type_glsl.svg",
	".frag":  "file_type_glsl.svg",
	".metal": "file_type_metal.svg",
	".wgsl":  "file_type_wgsl.svg",
	".slang": "file_type_slang.svg",

	// CUDA
	".cu":  "file_type_cuda.svg",
	".cuh": "file_type_cuda.svg",

	// Makefile / CMake / Build
	".cmake":   "file_type_cmake.svg",
	"Makefile": "file_type_gnu.svg",
	".mk":      "file_type_gnu.svg",

	// Diff / Patch
	".diff":  "file_type_diff.svg",
	".patch": "file_type_patch.svg",

	// Lock files
	".lock": "file_type_config.svg",

	// Images
	".png":  "file_type_image.svg",
	".jpg":  "file_type_image.svg",
	".jpeg": "file_type_image.svg",
	".gif":  "file_type_image.svg",
	".svg":  "file_type_svg.svg",
	".webp": "file_type_image.svg",
	".ico":  "file_type_favicon.svg",
	".avif": "file_type_avif.svg",

	// Fonts
	".ttf":   "file_type_font.svg",
	".otf":   "file_type_font.svg",
	".woff":  "file_type_font.svg",
	".woff2": "file_type_font.svg",
	".eot":   "file_type_font.svg",

	// Audio
	".mp3":  "file_type_audio.svg",
	".wav":  "file_type_audio.svg",
	".ogg":  "file_type_audio.svg",
	".flac": "file_type_audio.svg",

	// Video
	".mp4":  "file_type_video.svg",
	".webm": "file_type_video.svg",
	".mkv":  "file_type_video.svg",
	".mov":  "file_type_video.svg",
	".avi":  "file_type_video.svg",

	// Archives
	".zip": "file_type_zip.svg",
	".tar": "file_type_zip.svg",
	".gz":  "file_type_zip.svg",
	".bz2": "file_type_zip.svg",
	".rar": "file_type_zip.svg",
	".7z":  "file_type_zip.svg",

	// Binary
	".bin":   "file_type_binary.svg",
	".exe":   "file_type_binary.svg",
	".dll":   "file_type_binary.svg",
	".so":    "file_type_binary.svg",
	".dylib": "file_type_binary.svg",

	// Certificates / Keys
	".pem": "file_type_cert.svg",
	".crt": "file_type_cert.svg",
	".cer": "file_type_cert.svg",
	".key": "file_type_key.svg",
	".gpg": "file_type_gpg.svg",

	// Logs
	".log": "file_type_log.svg",

	// Backup
	".bak": "file_type_bak.svg",
}

// LanguageIconMap maps language names (as typically reported by editors/LSPs)
// to their icon SVG filename.
var LanguageIconMap = map[string]string{
	"go":              "file_type_go.svg",
	"javascript":      "file_type_js.svg",
	"javascriptreact": "file_type_reactjs.svg",
	"typescript":      "file_type_typescript.svg",
	"typescriptreact": "file_type_reactts.svg",
	"python":          "file_type_python.svg",
	"rust":            "file_type_rust.svg",
	"c":               "file_type_c.svg",
	"cpp":             "file_type_cpp.svg",
	"csharp":          "file_type_csharp.svg",
	"java":            "file_type_java.svg",
	"kotlin":          "file_type_kotlin.svg",
	"scala":           "file_type_scala.svg",
	"groovy":          "file_type_groovy.svg",
	"ruby":            "file_type_ruby.svg",
	"php":             "file_type_php.svg",
	"html":            "file_type_html.svg",
	"css":             "file_type_css.svg",
	"scss":            "file_type_scss.svg",
	"sass":            "file_type_sass.svg",
	"less":            "file_type_less.svg",
	"vue":             "file_type_vue.svg",
	"svelte":          "file_type_svelte.svg",
	"astro":           "file_type_astro.svg",
	"dart":            "file_type_dartlang.svg",
	"swift":           "file_type_swift.svg",
	"objective-c":     "file_type_objectivec.svg",
	"shellscript":     "file_type_shell.svg",
	"bash":            "file_type_shell.svg",
	"powershell":      "file_type_powershell.svg",
	"json":            "file_type_json.svg",
	"yaml":            "file_type_yaml.svg",
	"toml":            "file_type_toml.svg",
	"xml":             "file_type_xml.svg",
	"markdown":        "file_type_markdown.svg",
	"sql":             "file_type_sql.svg",
	"graphql":         "file_type_graphql.svg",
	"lua":             "file_type_lua.svg",
	"elixir":          "file_type_elixir.svg",
	"erlang":          "file_type_erlang.svg",
	"haskell":         "file_type_haskell.svg",
	"purescript":      "file_type_purescript.svg",
	"ocaml":           "file_type_ocaml.svg",
	"fsharp":          "file_type_fsharp.svg",
	"terraform":       "file_type_terraform.svg",
	"protobuf":        "file_type_protobuf.svg",
	"wasm":            "file_type_wasm.svg",
	"asm":             "file_type_assembly.svg",
	"r":               "file_type_r.svg",
	"julia":           "file_type_julia.svg",
	"nim":             "file_type_nim.svg",
	"zig":             "file_type_zig.svg",
	"v":               "file_type_vlang.svg",
	"clojure":         "file_type_clojure.svg",
	"clojurescript":   "file_type_clojurescript.svg",
	"solidity":        "file_type_solidity.svg",
	"crystal":         "file_type_crystal.svg",
	"perl":            "file_type_perl.svg",
	"coffeescript":    "file_type_coffeescript.svg",
	"d":               "file_type_dlang.svg",
	"mojo":            "file_type_mojo.svg",
	"gleam":           "file_type_gleam.svg",
	"nix":             "file_type_nix.svg",
	"idris":           "file_type_idris.svg",
	"agda":            "file_type_agda.svg",
	"lean":            "file_type_lean.svg",
	"gdscript":        "file_type_gdscript.svg",
	"glsl":            "file_type_glsl.svg",
	"hlsl":            "file_type_hlsl.svg",
	"cuda":            "file_type_cuda.svg",
	"dockerfile":      "file_type_docker.svg",
	"diff":            "file_type_diff.svg",
	"plaintext":       "file_type_text.svg",
	"log":             "file_type_log.svg",
	"batch":           "file_type_bat.svg",
	"matlab":          "file_type_matlab.svg",
	"fortran":         "file_type_fortran.svg",
	"cobol":           "file_type_cobol.svg",
	"ada":             "file_type_ada.svg",
	"prolog":          "file_type_prolog.svg",
	"tcl":             "file_type_tcl.svg",
	"vala":            "file_type_vala.svg",
	"wgsl":            "file_type_wgsl.svg",
}

// FrameworkIconMap maps framework/tool names to their icon SVG filename.
var FrameworkIconMap = map[string]string{
	// JavaScript / TypeScript frameworks
	"react":     "file_type_reactjs.svg",
	"next":      "file_type_next.svg",
	"nextjs":    "file_type_next.svg",
	"nuxt":      "file_type_nuxt.svg",
	"nuxtjs":    "file_type_nuxt.svg",
	"vue":       "file_type_vue.svg",
	"svelte":    "file_type_svelte.svg",
	"sveltekit": "file_type_svelte.svg",
	"astro":     "file_type_astro.svg",
	"angular":   "file_type_angular.svg",
	"ember":     "file_type_ember.svg",
	"gatsby":    "file_type_gatsby.svg",
	"remix":     "file_type_reactjs.svg",
	"preact":    "file_type_preact.svg",
	"qwik":      "file_type_js.svg",

	// Backend / Server
	"nestjs":    "file_type_nestjs.svg",
	"express":   "file_type_node.svg",
	"fastify":   "file_type_node.svg",
	"koa":       "file_type_node.svg",
	"hono":      "file_type_js.svg",
	"django":    "file_type_django.svg",
	"flask":     "file_type_python.svg",
	"fastapi":   "file_type_python.svg",
	"rails":     "file_type_rails.svg",
	"laravel":   "file_type_php.svg",
	"symfony":   "file_type_symfony.svg",
	"spring":    "file_type_java.svg",
	"quarkus":   "file_type_java.svg",
	"micronaut": "file_type_java.svg",
	"ktor":      "file_type_kotlin.svg",
	"gin":       "file_type_go.svg",
	"echo":      "file_type_go.svg",
	"fiber":     "file_type_go.svg",
	"actix":     "file_type_rust.svg",
	"axum":      "file_type_rust.svg",
	"phoenix":   "file_type_elixir.svg",

	// Mobile
	"flutter":      "file_type_flutter.svg",
	"react-native": "file_type_reactjs.svg",
	"expo":         "file_type_expo.svg",
	"ionic":        "file_type_ionic.svg",
	"capacitor":    "file_type_capacitor.svg",

	// Build / Tooling
	"vite":      "file_type_vite.svg",
	"webpack":   "file_type_webpack.svg",
	"rollup":    "file_type_rollup.svg",
	"esbuild":   "file_type_esbuild.svg",
	"parcel":    "file_type_js.svg",
	"babel":     "file_type_babel.svg",
	"turborepo": "file_type_turbo.svg",
	"nx":        "file_type_nx.svg",
	"gradle":    "file_type_gradle.svg",
	"maven":     "file_type_maven.svg",
	"cargo":     "file_type_cargo.svg",

	// Testing
	"jest":       "file_type_jest.svg",
	"vitest":     "file_type_vitest.svg",
	"cypress":    "file_type_cypress.svg",
	"playwright": "file_type_playwright.svg",
	"mocha":      "file_type_mocha.svg",
	"jasmine":    "file_type_jasmine.svg",
	"pytest":     "file_type_pytest.svg",

	// CSS Frameworks / Tools
	"tailwind": "file_type_tailwind.svg",
	"styled":   "file_type_styled.svg",
	"sass":     "file_type_sass.svg",
	"postcss":  "file_type_postcss.svg",
	"unocss":   "file_type_unocss.svg",

	// ORM / DB
	"prisma":    "file_type_prisma.svg",
	"mongoose":  "file_type_mongo.svg",
	"sequelize": "file_type_sequelize.svg",
	"typeorm":   "file_type_typescript.svg",
	"drizzle":   "file_type_sql.svg",

	// DevOps / Cloud
	"docker":         "file_type_docker.svg",
	"kubernetes":     "file_type_config.svg",
	"terraform":      "file_type_terraform.svg",
	"ansible":        "file_type_ansible.svg",
	"github-actions": "file_type_github.svg",
	"circleci":       "file_type_circleci.svg",
	"jenkins":        "file_type_jenkins.svg",
	"aws":            "file_type_aws.svg",
	"azure":          "file_type_azure.svg",
	"gcloud":         "file_type_gcloud.svg",

	// Game Engines
	"godot": "file_type_godot.svg",
	"unity": "file_type_shaderlab.svg",

	// Desktop
	"electron": "file_type_electron.svg",
	"tauri":    "file_type_tauri.svg",
}

// FolderIconMap maps common folder names to their icon SVG filename.
var FolderIconMap = map[string]string{
	"src":           "folder_type_src.svg",
	"source":        "folder_type_src.svg",
	"dist":          "folder_type_dist.svg",
	"build":         "folder_type_dist.svg",
	"out":           "folder_type_dist.svg",
	"public":        "folder_type_public.svg",
	"static":        "folder_type_public.svg",
	"assets":        "folder_type_asset.svg",
	"images":        "folder_type_images.svg",
	"img":           "folder_type_images.svg",
	"fonts":         "folder_type_fonts.svg",
	"audio":         "folder_type_audio.svg",
	"video":         "folder_type_video.svg",
	"css":           "folder_type_css.svg",
	"styles":        "folder_type_style.svg",
	"sass":          "folder_type_sass.svg",
	"scss":          "folder_type_sass.svg",
	"js":            "folder_type_js.svg",
	"scripts":       "folder_type_script.svg",
	"ts":            "folder_type_typescript.svg",
	"types":         "folder_type_typings.svg",
	"typings":       "folder_type_typings.svg",
	"components":    "folder_type_component.svg",
	"views":         "folder_type_view.svg",
	"pages":         "folder_type_view.svg",
	"layouts":       "folder_type_view.svg",
	"hooks":         "folder_type_hook.svg",
	"routes":        "folder_type_route.svg",
	"router":        "folder_type_route.svg",
	"api":           "folder_type_api.svg",
	"controllers":   "folder_type_controller.svg",
	"models":        "folder_type_model.svg",
	"services":      "folder_type_services.svg",
	"middleware":    "folder_type_middleware.svg",
	"helpers":       "folder_type_helper.svg",
	"utils":         "folder_type_tools.svg",
	"lib":           "folder_type_library.svg",
	"libs":          "folder_type_library.svg",
	"vendor":        "folder_type_library.svg",
	"node_modules":  "folder_type_node.svg",
	"packages":      "folder_type_package.svg",
	"pkg":           "folder_type_package.svg",
	"plugins":       "folder_type_plugin.svg",
	"modules":       "folder_type_module.svg",
	"config":        "folder_type_config.svg",
	"conf":          "folder_type_config.svg",
	"settings":      "folder_type_config.svg",
	"env":           "folder_type_config.svg",
	"test":          "folder_type_test.svg",
	"tests":         "folder_type_test.svg",
	"__tests__":     "folder_type_test.svg",
	"spec":          "folder_type_test.svg",
	"e2e":           "folder_type_e2e.svg",
	"coverage":      "folder_type_coverage.svg",
	"docs":          "folder_type_docs.svg",
	"documentation": "folder_type_docs.svg",
	"logs":          "folder_type_log.svg",
	"log":           "folder_type_log.svg",
	"tmp":           "folder_type_temp.svg",
	"temp":          "folder_type_temp.svg",
	"cache":         "folder_type_temp.svg",
	"db":            "folder_type_db.svg",
	"database":      "folder_type_db.svg",
	"migrations":    "folder_type_db.svg",
	"prisma":        "folder_type_prisma.svg",
	"graphql":       "folder_type_graphql.svg",
	"docker":        "folder_type_docker.svg",
	".github":       "folder_type_github.svg",
	".gitlab":       "folder_type_gitlab.svg",
	".git":          "folder_type_git.svg",
	".vscode":       "folder_type_vscode.svg",
	"android":       "folder_type_android.svg",
	"ios":           "folder_type_ios.svg",
	"python":        "folder_type_python.svg",
	"kotlin":        "folder_type_kotlin.svg",
	"dart":          "folder_type_dart.svg",
	"flutter":       "folder_type_flutter.svg",
	"redux":         "folder_type_redux.svg",
	"store":         "folder_type_redux.svg",
	"shared":        "folder_type_shared.svg",
	"common":        "folder_type_common.svg",
	"server":        "folder_type_server.svg",
	"client":        "folder_type_client.svg",
	"app":           "folder_type_app.svg",
	"svelte":        "folder_type_svelte.svg",
	"nuxt":          "folder_type_nuxt.svg",
	"next":          "folder_type_next.svg",
	"expo":          "folder_type_expo.svg",
	"internals":     "folder_type_src.svg",
	"internal":      "folder_type_src.svg",
	"cmd":           "folder_type_cli.svg",
	"cli":           "folder_type_cli.svg",
}

const defaultFileIcon = "default_file.svg"
const defaultFolderIcon = "default_folder.svg"

// GetFileIcon returns the best-matching icon SVG filename for a given file name.
// It first checks the full filename (e.g. "Dockerfile"), then the extension.
// Falls back to defaultFileIcon if nothing matches.
func GetFileIcon(filename string) string {
	// Check full filename first (e.g. "Dockerfile", "Makefile")
	if icon, ok := ExtensionIconMap[filename]; ok {
		return icon
	}
	// Walk extensions from longest to shortest to handle cases like ".d.ts"
	name := filename
	for {
		idx := -1
		for i, ch := range name {
			if ch == '.' && i > 0 {
				idx = i
				break
			}
		}
		if idx == -1 {
			break
		}
		ext := name[idx:]
		if icon, ok := ExtensionIconMap[ext]; ok {
			return icon
		}
		name = name[idx+1:]
	}
	return defaultFileIcon
}

// GetFolderIcon returns the best-matching icon SVG filename for a given folder name.
func GetFolderIcon(folderName string) string {
	if icon, ok := FolderIconMap[folderName]; ok {
		return icon
	}
	return defaultFolderIcon
}

// GetLanguageIcon returns the icon for a language identifier (lowercase).
func GetLanguageIcon(language string) string {
	if icon, ok := LanguageIconMap[language]; ok {
		return icon
	}
	return defaultFileIcon
}

// GetFrameworkIcon returns the icon for a framework/tool name (lowercase).
func GetFrameworkIcon(framework string) string {
	if icon, ok := FrameworkIconMap[framework]; ok {
		return icon
	}
	return defaultFileIcon
}
