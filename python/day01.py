from collections import defaultdict

from utils import read


def part1() -> None:
    arr1, arr2 = [], []
    with read(day=1, sample=False) as f:
        for line in f.readlines():
            p1, p2 = map(int, line.split())
            arr1.append(p1)
            arr2.append(p2)
        arr1 = sorted(arr1)
        arr2 = sorted(arr2)

    s = sum(abs(p1 - p2) for p1, p2 in zip(arr1, arr2, strict=False))
    print(s)


def part2() -> None:
    arr1 = []
    dict2: dict[int, int] = defaultdict(lambda: 0)
    with read(day=1, sample=False) as f:
        for line in f.readlines():
            p1, p2 = map(int, line.split())
            arr1.append(p1)
            dict2[p2] += 1
        arr1 = sorted(arr1)

    s = sum(p * dict2[p] for p in arr1)
    print(s)


if __name__ == "__main__":
    part1()
    part2()
