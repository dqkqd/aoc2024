[group('go')]
@go day:
  cd go && go run ./day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}

[group('python')]
@python day:
  cd python && uv run day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}.py

[group('rust')]
@rust day:
  cd rust && cargo run --bin day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}

