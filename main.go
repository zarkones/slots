package main

/*
	Don't mind the code. :)
	It has been made in an hour and represents a PoC.
*/

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"slot/assets"
	"slot/core"
	"strconv"
	"sync"
	"time"

	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func getID(img *image.Image) string {
	return fmt.Sprint((*img).At(50, 50))
}

var symbols = map[string]int{}

func main() {
	slot1Img, _ := png.Decode(bytes.NewReader(assets.Slot1))
	slot2Img, _ := png.Decode(bytes.NewReader(assets.Slot2))
	slot3Img, _ := png.Decode(bytes.NewReader(assets.Slot3))

	mainCont := container.NewWithoutLayout()
	mainCont.Resize(core.WindowSize)

	maxRows := float32(5)
	maxCols := float32(4)

	w := mainCont.Size().Width / maxRows
	h := mainCont.Size().Height / maxCols

	totalScore := 100
	spinCost := 20
	rewardTxt := widget.NewLabel(strconv.Itoa(totalScore))

	imgIndex := 0
	for i := float32(0); i < maxRows; i++ {
		for x := float32(0); x < maxCols; x++ {
			var img image.Image
			switch imgIndex {
			default:
				imgIndex = 0
				img = slot1Img
			case 1:
				img = slot2Img
			case 2:
				img = slot3Img

			}
			imgIndex++
			slot := canvas.NewImageFromImage(img)
			slot.Resize(fyne.NewSize(w, h))
			slot.Move(fyne.NewPos(i*w, x*h))
			mainCont.Add(slot)

			id := getID(&slot.Image)
			symbols[id] = rand.Intn(40) + 10
		}
	}

	spinDurationSecs := int64(1)

	getCol := func(col int) (slots []*fyne.CanvasObject) {
		slots = make([]*fyne.CanvasObject, int(maxCols))

		switch col {

		default:
			slots = make([]*fyne.CanvasObject, 0)
			return slots

		case 0:
			slots[0] = &mainCont.Objects[0]
			slots[1] = &mainCont.Objects[1]
			slots[2] = &mainCont.Objects[2]
			slots[3] = &mainCont.Objects[3]

		case 1:
			slots[0] = &mainCont.Objects[4]
			slots[1] = &mainCont.Objects[5]
			slots[2] = &mainCont.Objects[6]
			slots[3] = &mainCont.Objects[7]

		case 2:
			slots[0] = &mainCont.Objects[8]
			slots[1] = &mainCont.Objects[9]
			slots[2] = &mainCont.Objects[10]
			slots[3] = &mainCont.Objects[11]

		case 3:
			slots[0] = &mainCont.Objects[12]
			slots[1] = &mainCont.Objects[13]
			slots[2] = &mainCont.Objects[14]
			slots[3] = &mainCont.Objects[15]

		case 4:
			slots[0] = &mainCont.Objects[16]
			slots[1] = &mainCont.Objects[17]
			slots[2] = &mainCont.Objects[18]
			slots[3] = &mainCont.Objects[19]

		}

		return slots
	}

	lastSpin := time.Time{}
	spin := func() {
		now := time.Now()
		if now.Unix()-lastSpin.Unix() < spinDurationSecs {
			return
		}
		lastSpin = now

		fmt.Println("Spin")

		totalScore -= spinCost
		rewardTxt.SetText(strconv.Itoa(totalScore))

		speed := time.Millisecond * 80

		wg := &sync.WaitGroup{}

		for i := 0; i < int(maxRows); i++ {
			slots := getCol(i)
			wg.Add(1)

			go func() {
				defer wg.Done()

				for x := 0; x < rand.Intn(100)+10; x++ {

					for slotIndex, slotPtr := range slots {
						slot := *slotPtr

						if slot.Position().Y == 0 {
							trailingSlotIndex := slotIndex - 1
							if slotIndex == 0 {
								trailingSlotIndex = 3
							}
							trailingSlot := *slots[trailingSlotIndex]
							slot.Move(slot.Position().AddXY(0, trailingSlot.Position().Y+h))
						}

					}

					for _, slotPtr := range slots {
						slot := *slotPtr

						anim := canvas.NewPositionAnimation(slot.Position(), slot.Position().AddXY(0, -h), speed, func(p fyne.Position) {
							slot.Move(p)
						})

						anim.Start()
					}

					time.Sleep(speed + time.Millisecond*5)
				}

			}()

		}

		wg.Wait()

		scoreMap := map[int]int{}

		for _, slot := range mainCont.Objects {
			if slot.Position().Y != h*2 {
				continue
			}

			x := slot.(*canvas.Image)

			amount := symbols[getID(&x.Image)]

			if _, ok := scoreMap[amount]; !ok {
				scoreMap[amount] = 1
				continue
			}
			scoreMap[amount]++
		}

		for reward, f := range scoreMap {
			if f < 3 {
				continue
			}
			totalScore += reward * f
			rewardTxt.SetText(strconv.Itoa(totalScore))
		}

		if totalScore <= 0 {
			rewardTxt.SetText("Game Over")
		}
	}

	for i := 0; i < int(maxCols); i++ {
		line := canvas.NewRectangle(color.RGBA{255, 255, 255, 200})
		line.Resize(fyne.NewSize(w/40, mainCont.Size().Height))
		line.Move(fyne.NewPos(float32(1+i)*w-line.Size().Width/2, 0))
		mainCont.Add(line)
	}

	bg := canvas.NewRectangle(color.RGBA{0, 0, 0, 255})
	bg.Resize(fyne.NewSize(mainCont.Size().Width, h/3))
	rewardTxt.Move(fyne.NewPos(bg.Size().Width/2, bg.Size().Height/2))

	selector := canvas.NewRectangle(color.RGBA{10, 50, 255, 200})
	selector.Resize(fyne.NewSize(mainCont.Size().Width, h/20))
	selector.Move(fyne.NewPos(0, h*2+h/2-selector.Size().Height/2))

	core.MainW.SetContent(
		container.NewWithoutLayout(
			selector,
			mainCont,
			bg,
			rewardTxt,
		),
	)

	core.MainW.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		switch k.Name {
		case fyne.KeySpace:
			spin()
		}
	})

	core.MainW.ShowAndRun()
}
