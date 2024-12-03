import dataclasses
import io
from collections.abc import Callable, Generator

from utils import read


def peek(buf: io.TextIOWrapper) -> str:
    pos = buf.tell()
    c = buf.read(1)
    buf.seek(pos)
    return c


def readable(buf: io.TextIOWrapper) -> bool:
    return peek(buf) != ""


def read_while(
    buf: io.TextIOWrapper,
    pred: Callable[[str], bool],
) -> Generator[str]:
    while (c := peek(buf)) != "":
        if not pred(c):
            break
        buf.read(1)
        yield c


def expect(buf: io.TextIOWrapper, s: str) -> bool:
    pos = buf.tell()
    if buf.read(len(s)) == s:
        return True
    buf.seek(pos)
    return False


def read_int(buf: io.TextIOWrapper) -> int | None:
    num = "".join(read_while(buf, lambda x: x.isdigit()))
    return int(num) if num else None


@dataclasses.dataclass
class Do: ...


@dataclasses.dataclass
class Dont: ...


@dataclasses.dataclass
class Mul:
    lhs: int
    rhs: int


def read_mul(buf: io.TextIOWrapper) -> Mul | None:
    # consume
    for _ in read_while(buf, lambda b: b != "m"):
        continue

    if (
        expect(buf, "mul(")
        and ((lhs := read_int(buf)) is not None)
        and expect(buf, ",")
        and ((rhs := read_int(buf)) is not None)
        and expect(buf, ")")
    ):
        return Mul(lhs, rhs)

    # skip 1 to avoid infinite loop
    buf.read(1)
    return None


def read_all(buf: io.TextIOWrapper) -> Do | Dont | Mul | None:
    # consume
    for _ in read_while(buf, lambda b: b not in ["m", "d"]):
        continue

    if expect(buf, "do()"):
        return Do()

    if expect(buf, "don't()"):
        return Dont()

    return read_mul(buf)


def read_instructions() -> list[Do | Dont | Mul]:
    instructions = []
    with read(day=3, sample=False) as buf:
        while readable(buf):
            if (ins := read_all(buf)) is not None:
                instructions.append(ins)
    return instructions


def part1():
    ans = 0
    instructions = read_instructions()

    for ins in instructions:
        if isinstance(ins, Mul):
            ans += ins.lhs * ins.rhs
    print(ans)


def part2():
    ans = 0
    instructions = read_instructions()

    latest = Do()
    for ins in instructions:
        if isinstance(ins, Mul) and isinstance(latest, Do):
            ans += ins.lhs * ins.rhs
        elif not isinstance(ins, Mul):
            latest = ins

    print(ans)


if __name__ == "__main__":
    part1()
    part2()
