import dataclasses
import enum
from collections.abc import Callable, Generator

import utils


class Op(enum.Enum):
    Add = 0
    Mul = 1
    Concat = 2

    def eval(self, lhs: int, rhs: int) -> int:
        match self:
            case Op.Add:
                return lhs + rhs
            case Op.Mul:
                return lhs * rhs
            case Op.Concat:
                return int(f"{lhs}{rhs}")


def binary_ops_generator(size: int) -> Callable[[int], Generator[Op]]:
    def generator(mask: int) -> Generator[Op]:
        for b in range(size):
            if mask & (1 << b) > 0:
                yield Op.Add
            else:
                yield Op.Mul

    return generator


def ternary_ops_generator(size: int) -> Callable[[int], Generator[Op]]:
    def generator(mask: int) -> Generator[Op]:
        for b in range(size):
            s = (mask // (3**b)) % 3
            if s == 0:
                yield Op.Add
            elif s == 1:
                yield Op.Mul
            else:
                yield Op.Concat

    return generator


@dataclasses.dataclass
class Equation:
    lhs: int
    rhs: list[int]

    @classmethod
    def from_str(cls, line: str) -> "Equation | None":
        line = line.strip()
        if len(line) == 0:
            return None
        lhs_str, rhs_str = line.split(":")
        lhs = int(lhs_str)
        rhs = list(map(int, rhs_str.strip().split(" ")))
        return Equation(lhs=lhs, rhs=rhs)

    def good(self, ops: Generator[Op]) -> bool:
        value = self.rhs[0]
        for rhs, op in zip(self.rhs[1:], ops, strict=True):
            value = op.eval(value, rhs)
            if value > self.lhs:
                return False
        return value == self.lhs

    def from_add_and_mul(self) -> bool:
        gen = binary_ops_generator(len(self.rhs) - 1)
        return any(self.good(gen(mask)) for mask in range(1 << (len(self.rhs) - 1)))

    def from_all_ops(self) -> bool:
        gen = ternary_ops_generator(len(self.rhs) - 1)
        return any(self.good(gen(mask)) for mask in range(3 ** (len(self.rhs) - 1)))


def read_equations() -> Generator[Equation]:
    with utils.read(day=7, sample=False) as f:
        for line in f.readlines():
            equation = Equation.from_str(line)
            if equation is not None:
                yield equation


def part1() -> None:
    ans = 0
    for equation in read_equations():
        if equation.from_add_and_mul():
            ans += equation.lhs
    print(ans)


def part2() -> None:
    ans = 0
    for equation in read_equations():
        if equation.from_all_ops():
            ans += equation.lhs
    print(ans)


if __name__ == "__main__":
    part1()
    part2()
