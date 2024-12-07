import contextlib
from collections.abc import Generator
from io import TextIOWrapper
from pathlib import Path

INPUT_PATH = Path(__file__).parent.parent / "input"


@contextlib.contextmanager
def read(*, day: int, sample: bool) -> Generator[TextIOWrapper]:
    folder = INPUT_PATH / f"day{day:02d}"
    file = "sample.txt" if sample else "input.txt"
    with (folder / file).open("r") as f:
        yield f
