package cockpit_stream

import (
	"encoding/binary"
	"math/rand"
	"time"
)

func toByteArray(i int32) (arr [4]byte) {
	binary.BigEndian.PutUint32(arr[0:4], uint32(i))
	return
}

func generateRandomSlice(size int) []uint8 {
	slice := make([]uint8, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = uint8(rand.Intn(255))
	}
	return slice
}

//func main() {
//	//bytes := generateRandomSlice(100000)
//
//	bytes := generateRandomSlice(10_000)
//	bytes2 := generateRandomSlice(10_000)
//	msk := make([]uint8, 10_000)
//	for i := 0; i < 10; i++ {
//		msk[i] = bytes2[i] ^ bytes[i]
//	}
//
//	fmt.Println("minblocksize", lz4.CompressBlockBound(10_000))
//	hashtable := make([]int, 1<<16)
//	dest := make([]byte, lz4.CompressBlockBound(10_000))
//	l, err := lz4.CompressBlock(bytes, dest, hashtable)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("length", l)
//
//	//compressionTest()
//	//compressed := generateRandomSlice(8_000)
//
//	//start := time.Now()
//	//for i := 0; i < 10_000; i++ {
//	//	nonAppendMethod(compressed)
//	//}
//	//end := time.Now().Sub(start)
//	//fmt.Printf("non-append: %fs\n", end.Seconds())
//
//	//start = time.Now()
//	//for i := 0; i < 10_000; i++ {
//	//	appendMethod2(compressed)
//	//}
//	//end = time.Now().Sub(start)
//	//fmt.Printf("non-append2: %fs\n", end.Seconds())
//	//
//	//start = time.Now()
//	//for i := 0; i < 10_000; i++ {
//	//	appendMethod(compressed)
//	//}
//	//end = time.Now().Sub(start)
//	//fmt.Printf("non-append: %fs\n", end.Seconds())
//
//	//const fps int = 60
//	//const screens int = 10
//	//const frameCount = fps * 1
//	//start := time.Now()
//	////sleep := time.Duration(1000/fps) * time.Millisecond
//	//for scr := 1; scr <= frameCount; scr++ {
//	//	startCap := time.Now()
//	//	screenshot.Capture(0, 0, 500, 500)
//	//	fmt.Printf("cap time %d \n", time.Now().Sub(startCap).Milliseconds())
//	//
//	//	//time.Sleep(sleep - 20)
//	//}
//	//elapsed := time.Now().Sub(start)
//	//fmt.Printf("%f fps\n", float64(frameCount)/elapsed.Seconds())
//	//executeMutexScreenCap()
//	//execCallbackScreenCap()
//	//executeChannelScreenCap()
//}

var sizeOffset int = 1
var bufferOffset int = 5

func nonAppendMethod(compressed []uint8) {
	msg := make([]byte, 5+len(compressed))
	messageType := PayloadTypeMask
	var size int32 = 523_232

	msg[0] = byte(messageType)
	sizeBytes := toByteArray(size)
	for i := 0; i < len(sizeBytes); i++ {
		msg[i+sizeOffset] = sizeBytes[i]
	}

	for i := 0; i < len(compressed); i++ {
		msg[i+bufferOffset] = compressed[i]
	}
}

func appendMethod(compressed []uint8) {
	messageType := PayloadTypeMask
	var size int32 = 523_232

	sizeBytes := toByteArray(size)
	msg := []byte{byte(messageType)}
	msg = append(msg, sizeBytes[:]...)
	//msg := []byte{byte(messageType), sizeBytes[0], sizeBytes[1], sizeBytes[2], sizeBytes[3]}
	msg = append(msg, compressed...)
}
func appendMethod2(compressed []uint8) {
	messageType := PayloadTypeMask
	var size int32 = 523_232

	sizeBytes := toByteArray(size)
	//msg := []byte{byte(messageType)}
	//msg = append(msg, sizeBytes[:]...)
	msg := []byte{byte(messageType), sizeBytes[0], sizeBytes[1], sizeBytes[2], sizeBytes[3]}
	msg = append(msg, compressed...)
}

//func compressionTest() {
//	//SizeRect := image.Rect(0, 0, 500, 500)
//	width := 500
//	height := 500
//
//	next, _ := open("screen_next.png")
//	prev, _ := open("screen_prev.png")
//	//client, _ := open("screen_client.png")
//
//	compressionBuffer := NewBuffer(width, height)
//
//	//////////////////
//	// START SERVER
//	//////////////////
//	payload := NewPayloadEncoder()
//	payload.
//		SetPosX(20).
//		SetPosY(20).
//		SetHeight(uint32(height)).
//		SetWidth(uint32(width)).
//		SetType(PayloadTypeMask)
//
//	start := time.Now()
//	// calculate mask for server prev/next to send to client
//	CalculateBitmask(prev, next, compressionBuffer)
//	// compress bitmask
//	size, err := CompressBuffer(compressionBuffer)
//	if err != nil {
//		log.Fatal(err)
//	}
//	payload.SetBytes(compressionBuffer)
//
//	fmt.Println(size)
//	serverElapsed := time.Now().Sub(start)
//
//	//////////////////
//	// START CLIENT
//	//////////////////
//	clientStart := time.Now()
//	DecodeCompressedMask(payload.Bytes)
//	////// apply mask to client image
//	//clientIncomingMask := image.NewRGBA(SizeRect)
//	//clientIncomingMask.Pix = compressionBuffer.Bytes
//	//clientXor := NewXorBitmask(width, height)
//	//clientXor.CalculateBitmask(client, clientIncomingMask)
//	//
//	//newClient := image.NewRGBA(SizeRect)
//	//newClient.Pix = clientXor.Buffer
//	clientElapsed := time.Now().Sub(clientStart)
//	elapsed := time.Now().Sub(start)
//	//
//	fmt.Printf("server: %dms\n", serverElapsed.Milliseconds())
//	fmt.Printf("client: %dms\n", clientElapsed.Milliseconds())
//	fmt.Printf("total: %dms\n", elapsed.Milliseconds())
//	//
//	//save(newClient, "screen_client_new.png")
//
//}

//func execCallbackScreenCap() {
//	const fps int = 30
//	const screens int = 10
//	const frameCount = fps * 1
//	SizeRect := image.Rect(0, 0, 1000, 500)
//
//	capturer := screen_manager.New(&SizeRect, 60)
//	for i := 0; i < screens; i++ {
//		screen := screen_manager.NewVirtualScreen(i + 1)
//		capturer.RegisterScreen(screen)
//	}
//	capturer.Run()
//
//	time.Sleep(1 * time.Second)
//	capturer.Stop()
//}

//func executeMutexScreenCap() {
//
//	b1 := image.Rect(0, 0, 50, 50)
//	b2 := image.Rect(50, 50, 100, 100)
//	//img1 := image.NewRGBA(b1)
//	//img2 := image.NewRGBA(b2)
//	//screenshot.CaptureDisplay(0);
//	SizeRect := b1.Union(b2)
//
//	SizeRect = image.Rect(0, 0, 100, 100)
//	const fps int = 30
//	const screens int = 1
//	const frameCount = fps * 1
//	capturer := New(&SizeRect, fps)
//
//	capturer.Run()
//	var wg sync.WaitGroup
//	//time.Sleep(2 * time.Second)
//	for scr := 1; scr <= screens; scr++ {
//		b := image.Rect(0, 0, 500, 500)
//		img := image.NewRGBA(b)
//		wg.Add(1)
//		var curIdx int64 = -1
//		go func(idx int, wg *sync.WaitGroup) {
//			defer wg.Done()
//			start := time.Now()
//			i := 0
//			for i < frameCount {
//				startFrame := time.Now()
//				if capturer.Index() > curIdx {
//					capIdx := capturer.GetScreen(img, &b)
//					curIdx = capIdx
//					i++
//					//fmt.Printf("[%d]--- COPY FRAME %d ---\n", idx, curIdx)
//					elapsed := time.Now().Sub(startFrame).Milliseconds()
//					time.Sleep(time.Duration((1000/fps)-int(elapsed)) * time.Millisecond)
//				}
//			}
//			elapsed := time.Now().Sub(start)
//			fmt.Printf("img[%d]: %f fps\n", idx, float64(frameCount)/elapsed.Seconds())
//		}(scr, &wg)
//
//	}
//	wg.Wait()
//	//go func() {
//	//	start := time.Now()
//	//	for i := 0; i < 120; i++ {
//	//		capturer.GetScreen(img1, &b1)
//	//		time.Sleep(17 * time.Millisecond)
//	//	}
//	//	elapsed := time.Now().Sub(start)
//	//	fmt.Printf("img1: %f fps", float64(120)/elapsed.Seconds())
//	//}()
//	//go func() {
//	//	start := time.Now()
//	//	for i := 0; i < 120; i++ {
//	//		capturer.GetScreen(img2, &b2)
//	//		time.Sleep(17 * time.Millisecond)
//	//	}
//	//	elapsed := time.Now().Sub(start)
//	//	fmt.Printf("img1: %f fps", float64(120)/elapsed.Seconds())
//	//}()
//
//	fmt.Println("Stopping")
//	capturer.Stop()
//	fmt.Println("Shutting down")
//}
