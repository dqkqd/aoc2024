[group('go')]
@go day:
  cd go && go run ./day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}

[group('python')]
@python day:
  cd python && uv run day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}.py

[group('python')]
@python-lint:
  cd python && uv run ruff check . --fix --unsafe-fixes
  cd python && uv run ruff format
  cd python && uv run mypy . --explicit-package-bases

[group('rust')]
@rust day:
  cd rust && cargo run --bin day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}

sample day:
  mkdir -p input/day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}
  wl-paste > input/day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}/sample.txt

input day:
  mkdir -p input/day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}
  wl-paste > input/day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}/input.txt
