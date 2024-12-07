[group('go')]
@go day:
  cd go && go run ./day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}

[group('go')]
@go-all:
  #!/usr/bin/env sh
  cd go
  for day in `ls -d day*`; do
    echo
    echo "=====$day====="
    go run ./$day
  done

[group('python')]
@python day:
  cd python && uv run day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}.py

[group('python')]
@python-all:
  #!/usr/bin/env sh
  cd python
  for day in `ls day*`; do
    echo
    echo "=====$day====="
    uv run $day
  done

[group('python')]
@python-lint:
  cd python && uv run ruff check . --fix --unsafe-fixes
  cd python && uv run ruff format
  cd python && uv run mypy . --explicit-package-bases

[group('rust')]
@rust day:
  cd rust && cargo run --bin day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}

[group('rust')]
@rust-all:
  #!/usr/bin/env sh
  cd rust
  for day in `ls src/bin/day*`; do
    echo
    echo "=====${day:8:5}====="
    cargo run --release --bin ${day:8:5}
  done

sample day:
  mkdir -p input/day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}
  wl-paste > input/day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}/sample.txt

input day:
  mkdir -p input/day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}
  wl-paste > input/day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}/input.txt

sample-show day:
  cat input/day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}/sample.txt

input-show day:
  cat input/day{{ replace_regex(trim(day), "^(\\d)$", "0$1") }}/sample.txt
