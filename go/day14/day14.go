package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"strconv"
	"strings"

	utils "example.com/aoc2024"
	"github.com/golang/freetype"
	"github.com/icza/mjpeg"
)

var (
	fontfile = "/usr/share/fonts/TTF/JetBrainsMonoNerdFontMono-Bold.ttf"
	size     = 12.0
)

const (
	WIDTH  = 101
	HEIGHT = 103
)

type Tuple struct {
	x int
	y int
}

type (
	Point     Tuple
	Veclocity Tuple
)

type Robot struct {
	pos Point
	vec Veclocity
}

type Robots []Robot

func RobotFromString(s string) Robot {
	v := strings.Split(s, " ")

	readNumber := func(s string) (int, int) {
		ss := strings.Split(s, "=")[1]
		numbers := strings.Split(ss, ",")

		a, err := strconv.Atoi(numbers[0])
		utils.CheckError(err)

		b, err := strconv.Atoi(numbers[1])

		utils.CheckError(err)

		return a, b
	}

	px, py := readNumber(v[0])
	vx, vy := readNumber(v[1])

	return Robot{Point{px, py}, Veclocity{vx, vy}}
}

func (r *Robot) Run(seconds int) {
	pos := &r.pos
	pos.x += r.vec.x * seconds
	pos.y += r.vec.y * seconds
	pos.x %= WIDTH
	pos.y %= HEIGHT

	if pos.x < 0 {
		pos.x += WIDTH
	}

	if pos.y < 0 {
		pos.y += HEIGHT
	}
}

func ReadInput() Robots {
	buf := utils.Read(14, false)
	robots := []Robot{}
	for {
		s, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			robots = append(robots, RobotFromString(s))
		}
	}

	return robots
}

func Part1() {
	robots := ReadInput()
	for i := range robots {
		robots[i].Run(100)
	}

	numberOfRobots := []int{0, 0, 0, 0}

	middleHeight := HEIGHT / 2
	middleWidth := WIDTH / 2

	for _, robot := range robots {
		if robot.pos.x == middleWidth || robot.pos.y == middleHeight {
			continue
		}

		x := 0
		y := 0
		if robot.pos.x > middleWidth {
			x = 1
		}
		if robot.pos.y > middleHeight {
			y = 1
		}

		numberOfRobots[x*2+y]++
	}

	s := 1
	for _, x := range numberOfRobots {
		s *= x
	}

	fmt.Println(s)
}

func EmptyImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{0, 0}, draw.Src)
	return img
}

func (rs Robots) Render(index int, c *freetype.Context) *image.RGBA {
	img := EmptyImage()

	for _, r := range rs {
		img.Set(r.pos.y, r.pos.x, image.White)
	}

	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.White)
	pt := freetype.Pt(10, 10+int(c.PointToFixed(size)>>6))
	_, err := c.DrawString(strconv.Itoa(index), pt)
	utils.CheckError(err)

	return img
}

func Part2() {
	fontBytes, err := os.ReadFile(fontfile)
	utils.CheckError(err)
	f, err := freetype.ParseFont(fontBytes)
	utils.CheckError(err)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(size)

	robots := ReadInput()

	aw, err := mjpeg.New("test.avi", WIDTH, HEIGHT, 24)
	utils.CheckError(err)

	start := 7789

	for i := range robots {
		robots[i].Run(start)
	}

	for iter := start + 1; iter <= 7791; iter++ {
		for i := range robots {
			robots[i].Run(1)
		}
		img := robots.Render(iter, c)
		imageName := fmt.Sprintf("image_%d.jpeg", iter)

		f, err := os.Create(imageName)
		utils.CheckError(err)
		jpeg.Encode(f, img, nil)

		buf := &bytes.Buffer{}
		utils.CheckError(jpeg.Encode(buf, img, nil))
		aw.AddFrame(buf.Bytes())
	}

	utils.CheckError(aw.Close())
}

func main() {
	Part1()
	Part2()
}
