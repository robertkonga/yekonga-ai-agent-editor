/**
 * fileIcons.ts
 * Maps file extensions / filenames / language IDs / folder names to SVG icon
 * filenames that are served from the Vite static `public/file-icons/` directory.
 *
 * Mirrors the logic in desktop-go/icons/icon_map.go so both sides are in sync.
 */

// ─── Base path ────────────────────────────────────────────────────────────────
// Icons are symlinked → frontend/public/file-icons → desktop-go/icons/icons/
const BASE = '/file-icons/'

// ─── Extension → icon filename ────────────────────────────────────────────────
const EXT_MAP: Record<string, string> = {
  // Go
  '.go':      'file_type_go.svg',
  '.gomod':   'file_type_go_package.svg',
  '.gosum':   'file_type_go_package.svg',
  '.work':    'file_type_go_work.svg',

  // JavaScript
  '.js':    'file_type_js.svg',
  '.mjs':   'file_type_js.svg',
  '.cjs':   'file_type_js.svg',
  '.jsx':   'file_type_reactjs.svg',

  // TypeScript
  '.ts':    'file_type_typescript.svg',
  '.tsx':   'file_type_reactts.svg',
  '.d.ts':  'file_type_typescriptdef.svg',

  // Python
  '.py':    'file_type_python.svg',
  '.pyw':   'file_type_python.svg',
  '.pyx':   'file_type_cython.svg',
  '.ipynb': 'file_type_jupyter.svg',

  // Rust
  '.rs': 'file_type_rust.svg',

  // C / C++
  '.c':   'file_type_c.svg',
  '.h':   'file_type_cheader.svg',
  '.cpp': 'file_type_cpp.svg',
  '.cc':  'file_type_cpp.svg',
  '.cxx': 'file_type_cpp.svg',
  '.hpp': 'file_type_cppheader.svg',
  '.hh':  'file_type_cppheader.svg',

  // C#
  '.cs':     'file_type_csharp.svg',
  '.csproj': 'file_type_csproj.svg',

  // Java / JVM
  '.java':   'file_type_java.svg',
  '.jar':    'file_type_jar.svg',
  '.class':  'file_type_class.svg',
  '.kt':     'file_type_kotlin.svg',
  '.kts':    'file_type_kotlin.svg',
  '.groovy': 'file_type_groovy.svg',
  '.gradle': 'file_type_gradle.svg',
  '.scala':  'file_type_scala.svg',
  '.sbt':    'file_type_sbt.svg',

  // Ruby
  '.rb':      'file_type_ruby.svg',
  '.erb':     'file_type_erb.svg',
  '.rbs':     'file_type_ruby.svg',
  '.gemspec': 'file_type_bundler.svg',

  // PHP
  '.php':        'file_type_php.svg',
  '.phtml':      'file_type_php.svg',
  '.blade.php':  'file_type_blade.svg',

  // HTML / Web
  '.html':  'file_type_html.svg',
  '.htm':   'file_type_html.svg',
  '.xhtml': 'file_type_html.svg',

  // CSS
  '.css':  'file_type_css.svg',
  '.scss': 'file_type_scss.svg',
  '.sass': 'file_type_sass.svg',
  '.less': 'file_type_less.svg',
  '.styl': 'file_type_stylus.svg',

  // SFC frameworks
  '.vue':    'file_type_vue.svg',
  '.svelte': 'file_type_svelte.svg',
  '.astro':  'file_type_astro.svg',

  // Dart / Flutter
  '.dart': 'file_type_dartlang.svg',

  // Swift
  '.swift': 'file_type_swift.svg',

  // Objective-C
  '.m':  'file_type_objectivec.svg',
  '.mm': 'file_type_objectivecpp.svg',

  // Shell / Scripts
  '.sh':   'file_type_shell.svg',
  '.bash': 'file_type_shell.svg',
  '.zsh':  'file_type_shell.svg',
  '.fish': 'file_type_shell.svg',
  '.bat':  'file_type_bat.svg',
  '.cmd':  'file_type_bat.svg',
  '.ps1':  'file_type_powershell.svg',
  '.psm1': 'file_type_powershell_psm.svg',
  '.psd1': 'file_type_powershell_psd.svg',

  // Data / Config
  '.json':  'file_type_json.svg',
  '.json5': 'file_type_json5.svg',
  '.jsonc': 'file_type_json.svg',
  '.yaml':  'file_type_yaml.svg',
  '.yml':   'file_type_yaml.svg',
  '.toml':  'file_type_toml.svg',
  '.ini':   'file_type_ini.svg',
  '.env':   'file_type_dotenv.svg',
  '.xml':   'file_type_xml.svg',
  '.xsl':   'file_type_xsl.svg',
  '.xsd':   'file_type_xml.svg',
  '.csv':   'file_type_excel.svg',

  // Markdown / Docs
  '.md':  'file_type_markdown.svg',
  '.mdx': 'file_type_mdx.svg',
  '.rst': 'file_type_rest.svg',
  '.txt': 'file_type_text.svg',
  '.tex': 'file_type_tex.svg',
  '.pdf': 'file_type_pdf.svg',

  // SQL / Databases
  '.sql':    'file_type_sql.svg',
  '.pgsql':  'file_type_pgsql.svg',
  '.mysql':  'file_type_mysql.svg',
  '.db':     'file_type_db.svg',
  '.sqlite': 'file_type_sqlite.svg',

  // GraphQL
  '.graphql': 'file_type_graphql.svg',
  '.gql':     'file_type_graphql.svg',

  // Lua
  '.lua': 'file_type_lua.svg',

  // Elixir / Erlang
  '.ex':   'file_type_elixir.svg',
  '.exs':  'file_type_elixir.svg',
  '.eex':  'file_type_eex.svg',
  '.heex': 'file_type_eex.svg',
  '.erl':  'file_type_erlang.svg',
  '.hrl':  'file_type_erlang.svg',

  // Haskell / PureScript
  '.hs':   'file_type_haskell.svg',
  '.lhs':  'file_type_haskell.svg',
  '.purs': 'file_type_purescript.svg',

  // OCaml
  '.ml':  'file_type_ocaml.svg',
  '.mli': 'file_type_ocaml_intf.svg',

  // F#
  '.fs':     'file_type_fsharp.svg',
  '.fsi':    'file_type_fsharp.svg',
  '.fsx':    'file_type_fsharp.svg',
  '.fsproj': 'file_type_fsproj.svg',

  // Terraform / HCL
  '.tf':     'file_type_terraform.svg',
  '.tfvars': 'file_type_terraform.svg',
  '.hcl':    'file_type_hashicorp.svg',

  // Protobuf
  '.proto': 'file_type_protobuf.svg',

  // WASM
  '.wasm': 'file_type_wasm.svg',
  '.wat':  'file_type_wasm.svg',

  // Assembly
  '.asm':  'file_type_assembly.svg',
  '.s':    'file_type_assembly.svg',
  '.nasm': 'file_type_assembly.svg',

  // R
  '.r':    'file_type_r.svg',
  '.R':    'file_type_r.svg',
  '.rmd':  'file_type_rmd.svg',

  // Julia
  '.jl': 'file_type_julia.svg',

  // Nim
  '.nim':    'file_type_nim.svg',
  '.nimble': 'file_type_nimble.svg',

  // Zig
  '.zig': 'file_type_zig.svg',

  // V (Vlang)
  '.v': 'file_type_vlang.svg',

  // Clojure
  '.clj':  'file_type_clojure.svg',
  '.cljs': 'file_type_clojurescript.svg',
  '.cljc': 'file_type_clojure.svg',
  '.edn':  'file_type_clojure.svg',

  // Solidity
  '.sol': 'file_type_solidity.svg',

  // Crystal
  '.cr': 'file_type_crystal.svg',

  // Perl
  '.pl':  'file_type_perl.svg',
  '.pm':  'file_type_perl.svg',
  '.pod': 'file_type_perl.svg',

  // Gleam
  '.gleam': 'file_type_gleam.svg',

  // Nix
  '.nix': 'file_type_nix.svg',

  // Idris / Agda / Lean
  '.idr':  'file_type_idris.svg',
  '.agda': 'file_type_agda.svg',
  '.lean': 'file_type_lean.svg',

  // GDScript
  '.gd': 'file_type_gdscript.svg',

  // Shaders
  '.hlsl':  'file_type_hlsl.svg',
  '.glsl':  'file_type_glsl.svg',
  '.vert':  'file_type_glsl.svg',
  '.frag':  'file_type_glsl.svg',
  '.metal': 'file_type_metal.svg',
  '.wgsl':  'file_type_wgsl.svg',
  '.slang': 'file_type_slang.svg',

  // CUDA
  '.cu':  'file_type_cuda.svg',
  '.cuh': 'file_type_cuda.svg',

  // Build
  '.cmake':  'file_type_cmake.svg',
  '.mk':     'file_type_gnu.svg',

  // Diff / Patch
  '.diff':  'file_type_diff.svg',
  '.patch': 'file_type_patch.svg',

  // Lock
  '.lock': 'file_type_config.svg',

  // Images
  '.png':  'file_type_image.svg',
  '.jpg':  'file_type_image.svg',
  '.jpeg': 'file_type_image.svg',
  '.gif':  'file_type_image.svg',
  '.svg':  'file_type_svg.svg',
  '.webp': 'file_type_image.svg',
  '.ico':  'file_type_favicon.svg',
  '.avif': 'file_type_avif.svg',

  // Fonts
  '.ttf':   'file_type_font.svg',
  '.otf':   'file_type_font.svg',
  '.woff':  'file_type_font.svg',
  '.woff2': 'file_type_font.svg',

  // Audio
  '.mp3':  'file_type_audio.svg',
  '.wav':  'file_type_audio.svg',
  '.ogg':  'file_type_audio.svg',
  '.flac': 'file_type_audio.svg',

  // Video
  '.mp4':  'file_type_video.svg',
  '.webm': 'file_type_video.svg',
  '.mkv':  'file_type_video.svg',
  '.mov':  'file_type_video.svg',

  // Archives
  '.zip': 'file_type_zip.svg',
  '.tar': 'file_type_zip.svg',
  '.gz':  'file_type_zip.svg',
  '.bz2': 'file_type_zip.svg',
  '.rar': 'file_type_zip.svg',
  '.7z':  'file_type_zip.svg',

  // Binary
  '.bin':   'file_type_binary.svg',
  '.exe':   'file_type_binary.svg',
  '.dll':   'file_type_binary.svg',
  '.so':    'file_type_binary.svg',
  '.dylib': 'file_type_binary.svg',

  // Certs / Keys
  '.pem': 'file_type_cert.svg',
  '.crt': 'file_type_cert.svg',
  '.cer': 'file_type_cert.svg',
  '.key': 'file_type_key.svg',
  '.gpg': 'file_type_gpg.svg',

  // Log / Backup
  '.log': 'file_type_log.svg',
  '.bak': 'file_type_bak.svg',
}

// ─── Exact filename → icon filename ──────────────────────────────────────────
const FILENAME_MAP: Record<string, string> = {
  'Dockerfile':          'file_type_docker.svg',
  'docker-compose.yml':  'file_type_docker.svg',
  'docker-compose.yaml': 'file_type_docker.svg',
  '.dockerignore':       'file_type_docker.svg',
  'Makefile':            'file_type_gnu.svg',
  'Gemfile':             'file_type_bundler.svg',
  'Rakefile':            'file_type_rake.svg',
  '.gitignore':          'file_type_git.svg',
  '.gitattributes':      'file_type_git.svg',
  '.env':                'file_type_dotenv.svg',
  '.env.local':          'file_type_dotenv.svg',
  '.env.example':        'file_type_dotenv.svg',
  'go.mod':              'file_type_go_package.svg',
  'go.sum':              'file_type_go_package.svg',
  'package.json':        'file_type_npm.svg',
  'package-lock.json':   'file_type_npm.svg',
  'tsconfig.json':       'file_type_tsconfig.svg',
  'vite.config.ts':      'file_type_vite.svg',
  'vite.config.js':      'file_type_vite.svg',
  'webpack.config.js':   'file_type_webpack.svg',
  'rollup.config.js':    'file_type_rollup.svg',
  '.eslintrc':           'file_type_eslint.svg',
  '.eslintrc.js':        'file_type_eslint.svg',
  '.eslintrc.json':      'file_type_eslint.svg',
  '.prettierrc':         'file_type_prettier.svg',
  '.prettierrc.js':      'file_type_prettier.svg',
  '.babelrc':            'file_type_babel.svg',
  'babel.config.js':     'file_type_babel.svg',
  'tailwind.config.js':  'file_type_tailwind.svg',
  'tailwind.config.ts':  'file_type_tailwind.svg',
  'astro.config.mjs':    'file_type_astroconfig.svg',
  'svelte.config.js':    'file_type_svelteconfig.svg',
  'cargo.toml':          'file_type_cargo.svg',
  'Cargo.toml':          'file_type_cargo.svg',
}

// ─── Language ID → icon filename ─────────────────────────────────────────────
const LANG_MAP: Record<string, string> = {
  go:               'file_type_go.svg',
  javascript:       'file_type_js.svg',
  javascriptreact:  'file_type_reactjs.svg',
  typescript:       'file_type_typescript.svg',
  typescriptreact:  'file_type_reactts.svg',
  python:           'file_type_python.svg',
  rust:             'file_type_rust.svg',
  c:                'file_type_c.svg',
  cpp:              'file_type_cpp.svg',
  csharp:           'file_type_csharp.svg',
  java:             'file_type_java.svg',
  kotlin:           'file_type_kotlin.svg',
  scala:            'file_type_scala.svg',
  groovy:           'file_type_groovy.svg',
  ruby:             'file_type_ruby.svg',
  php:              'file_type_php.svg',
  html:             'file_type_html.svg',
  css:              'file_type_css.svg',
  scss:             'file_type_scss.svg',
  sass:             'file_type_sass.svg',
  less:             'file_type_less.svg',
  vue:              'file_type_vue.svg',
  svelte:           'file_type_svelte.svg',
  astro:            'file_type_astro.svg',
  dart:             'file_type_dartlang.svg',
  swift:            'file_type_swift.svg',
  'objective-c':    'file_type_objectivec.svg',
  shellscript:      'file_type_shell.svg',
  shell:            'file_type_shell.svg',
  bash:             'file_type_shell.svg',
  powershell:       'file_type_powershell.svg',
  json:             'file_type_json.svg',
  yaml:             'file_type_yaml.svg',
  toml:             'file_type_toml.svg',
  xml:              'file_type_xml.svg',
  markdown:         'file_type_markdown.svg',
  sql:              'file_type_sql.svg',
  graphql:          'file_type_graphql.svg',
  lua:              'file_type_lua.svg',
  elixir:           'file_type_elixir.svg',
  erlang:           'file_type_erlang.svg',
  haskell:          'file_type_haskell.svg',
  purescript:       'file_type_purescript.svg',
  ocaml:            'file_type_ocaml.svg',
  fsharp:           'file_type_fsharp.svg',
  terraform:        'file_type_terraform.svg',
  protobuf:         'file_type_protobuf.svg',
  wasm:             'file_type_wasm.svg',
  asm:              'file_type_assembly.svg',
  r:                'file_type_r.svg',
  julia:            'file_type_julia.svg',
  nim:              'file_type_nim.svg',
  zig:              'file_type_zig.svg',
  v:                'file_type_vlang.svg',
  clojure:          'file_type_clojure.svg',
  clojurescript:    'file_type_clojurescript.svg',
  solidity:         'file_type_solidity.svg',
  crystal:          'file_type_crystal.svg',
  perl:             'file_type_perl.svg',
  coffeescript:     'file_type_coffeescript.svg',
  d:                'file_type_dlang.svg',
  mojo:             'file_type_mojo.svg',
  gleam:            'file_type_gleam.svg',
  nix:              'file_type_nix.svg',
  idris:            'file_type_idris.svg',
  agda:             'file_type_agda.svg',
  lean:             'file_type_lean.svg',
  gdscript:         'file_type_gdscript.svg',
  glsl:             'file_type_glsl.svg',
  hlsl:             'file_type_hlsl.svg',
  cuda:             'file_type_cuda.svg',
  dockerfile:       'file_type_docker.svg',
  diff:             'file_type_diff.svg',
  plaintext:        'file_type_text.svg',
  log:              'file_type_log.svg',
  batch:            'file_type_bat.svg',
  matlab:           'file_type_matlab.svg',
  fortran:          'file_type_fortran.svg',
  cobol:            'file_type_cobol.svg',
  ada:              'file_type_ada.svg',
  prolog:           'file_type_prolog.svg',
  tcl:              'file_type_tcl.svg',
  vala:             'file_type_vala.svg',
  wgsl:             'file_type_wgsl.svg',
  ini:              'file_type_ini.svg',
}

// ─── Folder name → icon filenames ────────────────────────────────────────────
const FOLDER_MAP: Record<string, { closed: string; open: string }> = {
  src:           { closed: 'folder_type_src.svg',          open: 'folder_type_src_opened.svg' },
  source:        { closed: 'folder_type_src.svg',          open: 'folder_type_src_opened.svg' },
  dist:          { closed: 'folder_type_dist.svg',         open: 'folder_type_dist_opened.svg' },
  build:         { closed: 'folder_type_dist.svg',         open: 'folder_type_dist_opened.svg' },
  out:           { closed: 'folder_type_dist.svg',         open: 'folder_type_dist_opened.svg' },
  public:        { closed: 'folder_type_public.svg',       open: 'folder_type_public_opened.svg' },
  static:        { closed: 'folder_type_public.svg',       open: 'folder_type_public_opened.svg' },
  assets:        { closed: 'folder_type_asset.svg',        open: 'folder_type_asset_opened.svg' },
  images:        { closed: 'folder_type_images.svg',       open: 'folder_type_images_opened.svg' },
  img:           { closed: 'folder_type_images.svg',       open: 'folder_type_images_opened.svg' },
  fonts:         { closed: 'folder_type_fonts.svg',        open: 'folder_type_fonts_opened.svg' },
  audio:         { closed: 'folder_type_audio.svg',        open: 'folder_type_audio_opened.svg' },
  video:         { closed: 'folder_type_video.svg',        open: 'folder_type_video_opened.svg' },
  css:           { closed: 'folder_type_css.svg',          open: 'folder_type_css_opened.svg' },
  styles:        { closed: 'folder_type_style.svg',        open: 'folder_type_style_opened.svg' },
  sass:          { closed: 'folder_type_sass.svg',         open: 'folder_type_sass_opened.svg' },
  scss:          { closed: 'folder_type_sass.svg',         open: 'folder_type_sass_opened.svg' },
  js:            { closed: 'folder_type_js.svg',           open: 'folder_type_js_opened.svg' },
  scripts:       { closed: 'folder_type_script.svg',       open: 'folder_type_script_opened.svg' },
  ts:            { closed: 'folder_type_typescript.svg',   open: 'folder_type_typescript_opened.svg' },
  types:         { closed: 'folder_type_typings.svg',      open: 'folder_type_typings_opened.svg' },
  typings:       { closed: 'folder_type_typings.svg',      open: 'folder_type_typings_opened.svg' },
  components:    { closed: 'folder_type_component.svg',    open: 'folder_type_component_opened.svg' },
  views:         { closed: 'folder_type_view.svg',         open: 'folder_type_view_opened.svg' },
  pages:         { closed: 'folder_type_view.svg',         open: 'folder_type_view_opened.svg' },
  layouts:       { closed: 'folder_type_view.svg',         open: 'folder_type_view_opened.svg' },
  hooks:         { closed: 'folder_type_hook.svg',         open: 'folder_type_hook_opened.svg' },
  routes:        { closed: 'folder_type_route.svg',        open: 'folder_type_route_opened.svg' },
  router:        { closed: 'folder_type_route.svg',        open: 'folder_type_route_opened.svg' },
  api:           { closed: 'folder_type_api.svg',          open: 'folder_type_api_opened.svg' },
  controllers:   { closed: 'folder_type_controller.svg',   open: 'folder_type_controller_opened.svg' },
  models:        { closed: 'folder_type_model.svg',        open: 'folder_type_model_opened.svg' },
  services:      { closed: 'folder_type_services.svg',     open: 'folder_type_services_opened.svg' },
  middleware:    { closed: 'folder_type_middleware.svg',   open: 'folder_type_middleware_opened.svg' },
  helpers:       { closed: 'folder_type_helper.svg',       open: 'folder_type_helper_opened.svg' },
  utils:         { closed: 'folder_type_tools.svg',        open: 'folder_type_tools_opened.svg' },
  lib:           { closed: 'folder_type_library.svg',      open: 'folder_type_library_opened.svg' },
  libs:          { closed: 'folder_type_library.svg',      open: 'folder_type_library_opened.svg' },
  vendor:        { closed: 'folder_type_library.svg',      open: 'folder_type_library_opened.svg' },
  node_modules:  { closed: 'folder_type_node.svg',         open: 'folder_type_node_opened.svg' },
  packages:      { closed: 'folder_type_package.svg',      open: 'folder_type_package_opened.svg' },
  pkg:           { closed: 'folder_type_package.svg',      open: 'folder_type_package_opened.svg' },
  plugins:       { closed: 'folder_type_plugin.svg',       open: 'folder_type_plugin_opened.svg' },
  modules:       { closed: 'folder_type_module.svg',       open: 'folder_type_module_opened.svg' },
  config:        { closed: 'folder_type_config.svg',       open: 'folder_type_config_opened.svg' },
  conf:          { closed: 'folder_type_config.svg',       open: 'folder_type_config_opened.svg' },
  settings:      { closed: 'folder_type_config.svg',       open: 'folder_type_config_opened.svg' },
  env:           { closed: 'folder_type_config.svg',       open: 'folder_type_config_opened.svg' },
  test:          { closed: 'folder_type_test.svg',         open: 'folder_type_test_opened.svg' },
  tests:         { closed: 'folder_type_test.svg',         open: 'folder_type_test_opened.svg' },
  '__tests__':   { closed: 'folder_type_test.svg',         open: 'folder_type_test_opened.svg' },
  spec:          { closed: 'folder_type_test.svg',         open: 'folder_type_test_opened.svg' },
  e2e:           { closed: 'folder_type_e2e.svg',          open: 'folder_type_e2e_opened.svg' },
  coverage:      { closed: 'folder_type_coverage.svg',     open: 'folder_type_coverage_opened.svg' },
  docs:          { closed: 'folder_type_docs.svg',         open: 'folder_type_docs_opened.svg' },
  documentation: { closed: 'folder_type_docs.svg',         open: 'folder_type_docs_opened.svg' },
  logs:          { closed: 'folder_type_log.svg',          open: 'folder_type_log_opened.svg' },
  log:           { closed: 'folder_type_log.svg',          open: 'folder_type_log_opened.svg' },
  tmp:           { closed: 'folder_type_temp.svg',         open: 'folder_type_temp_opened.svg' },
  temp:          { closed: 'folder_type_temp.svg',         open: 'folder_type_temp_opened.svg' },
  cache:         { closed: 'folder_type_temp.svg',         open: 'folder_type_temp_opened.svg' },
  db:            { closed: 'folder_type_db.svg',           open: 'folder_type_db_opened.svg' },
  database:      { closed: 'folder_type_db.svg',           open: 'folder_type_db_opened.svg' },
  migrations:    { closed: 'folder_type_db.svg',           open: 'folder_type_db_opened.svg' },
  prisma:        { closed: 'folder_type_prisma.svg',       open: 'folder_type_prisma_opened.svg' },
  graphql:       { closed: 'folder_type_graphql.svg',      open: 'folder_type_graphql_opened.svg' },
  docker:        { closed: 'folder_type_docker.svg',       open: 'folder_type_docker_opened.svg' },
  '.github':     { closed: 'folder_type_github.svg',       open: 'folder_type_github_opened.svg' },
  '.gitlab':     { closed: 'folder_type_gitlab.svg',       open: 'folder_type_gitlab_opened.svg' },
  '.git':        { closed: 'folder_type_git.svg',          open: 'folder_type_git_opened.svg' },
  '.vscode':     { closed: 'folder_type_vscode.svg',       open: 'folder_type_vscode_opened.svg' },
  android:       { closed: 'folder_type_android.svg',      open: 'folder_type_android_opened.svg' },
  ios:           { closed: 'folder_type_ios.svg',          open: 'folder_type_ios_opened.svg' },
  python:        { closed: 'folder_type_python.svg',       open: 'folder_type_python_opened.svg' },
  kotlin:        { closed: 'folder_type_kotlin.svg',       open: 'folder_type_kotlin_opened.svg' },
  dart:          { closed: 'folder_type_dart.svg',         open: 'folder_type_dart_opened.svg' },
  flutter:       { closed: 'folder_type_flutter.svg',      open: 'folder_type_flutter_opened.svg' },
  redux:         { closed: 'folder_type_redux.svg',        open: 'folder_type_redux_opened.svg' },
  store:         { closed: 'folder_type_redux.svg',        open: 'folder_type_redux_opened.svg' },
  shared:        { closed: 'folder_type_shared.svg',       open: 'folder_type_shared_opened.svg' },
  common:        { closed: 'folder_type_common.svg',       open: 'folder_type_common_opened.svg' },
  server:        { closed: 'folder_type_server.svg',       open: 'folder_type_server_opened.svg' },
  client:        { closed: 'folder_type_client.svg',       open: 'folder_type_client_opened.svg' },
  app:           { closed: 'folder_type_app.svg',          open: 'folder_type_app_opened.svg' },
  svelte:        { closed: 'folder_type_svelte.svg',       open: 'folder_type_svelte_opened.svg' },
  nuxt:          { closed: 'folder_type_nuxt.svg',         open: 'folder_type_nuxt_opened.svg' },
  next:          { closed: 'folder_type_next.svg',         open: 'folder_type_next_opened.svg' },
  expo:          { closed: 'folder_type_expo.svg',         open: 'folder_type_expo_opened.svg' },
  internals:     { closed: 'folder_type_src.svg',          open: 'folder_type_src_opened.svg' },
  internal:      { closed: 'folder_type_src.svg',          open: 'folder_type_src_opened.svg' },
  cmd:           { closed: 'folder_type_cli.svg',          open: 'folder_type_cli_opened.svg' },
  cli:           { closed: 'folder_type_cli.svg',          open: 'folder_type_cli_opened.svg' },
}

const DEFAULT_FILE   = 'default_file.svg'
const DEFAULT_FOLDER = 'folder_type_src.svg'

// ─── Public API ───────────────────────────────────────────────────────────────

/**
 * Returns the full URL path to the SVG icon for a given file.
 * @param filename  Full filename (e.g. "main.go", "Dockerfile", ".env")
 * @param lang      Optional language ID from the store (e.g. "typescript")
 */
export function getFileIcon(filename: string, lang?: string): string {
  const lower = filename.trim().toLowerCase()

  // 1. Exact filename match (highest priority)
  if (FILENAME_MAP[filename]) return BASE + FILENAME_MAP[filename]
  if (FILENAME_MAP[lower])    return BASE + FILENAME_MAP[lower]

  // 2. Multi-part extension (e.g. ".d.ts", ".blade.php")
  const parts = lower.split('.')
  for (let i = 1; i < parts.length; i++) {
    const ext = '.' + parts.slice(i).join('.')
    if (EXT_MAP[ext]) return BASE + EXT_MAP[ext]
  }

  // 3. Language ID fallback
  if (lang && LANG_MAP[lang]) return BASE + LANG_MAP[lang]

  return BASE + DEFAULT_FILE
}

/**
 * Returns the full URL path to the SVG icon for a folder.
 * @param folderName  Folder name (e.g. "src", "node_modules", ".github")
 * @param expanded    Whether the folder is currently open/expanded
 */
export function getFolderIcon(folderName: string, expanded = false): string {
  const entry = FOLDER_MAP[folderName] ?? FOLDER_MAP[folderName.toLowerCase()]
  if (!entry) return BASE + (expanded ? 'default_folder_opened.svg' : DEFAULT_FOLDER)
  return BASE + (expanded ? entry.open : entry.closed)
}

/**
 * Returns the full URL path to the SVG icon by language ID alone.
 */
export function getLangIcon(lang: string): string {
  return BASE + (LANG_MAP[lang] ?? DEFAULT_FILE)
}
