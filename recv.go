package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/schollz/croc/v10/src/croc"
	log "github.com/schollz/logger"
)

func recvTabItem(a fyne.App, w fyne.Window) *container.TabItem {
	status := widget.NewLabel("")
	defer func() {
		if r := recover(); r != nil {
			log.Error(fmt.Sprint(r))
		}
	}()
	prog := widget.NewProgressBar()
	prog.Hide()

	topline := widget.NewLabel(lp(""))
	recvEntry := widget.NewEntry()
	recvEntry.SetPlaceHolder(lp("Enter code to download"))

	recvDir, _ := os.MkdirTemp("", "crocgui-recv")

	boxholder := container.NewVBox()
	receiverScroller := container.NewVScroll(boxholder)
	fileentries := make(map[string]*fyne.Container)

	var lastSaveDir string

	debugBox := container.NewHBox(widget.NewLabel(lp("Debug log:")), layout.NewSpacer(), widget.NewButton("Export full log", func() {
		savedialog := dialog.NewFileSave(func(f fyne.URIWriteCloser, e error) {
			if f != nil {
				logoutput.buf.WriteTo(f)
				f.Close()
			}
		}, w)
		savedialog.SetFileName("crocdebuglog.txt")
		savedialog.Show()
	}))
	debugObjects = append(debugObjects, debugBox)

	cancelchan := make(chan bool)
	activeButtonHolder := container.NewVBox()
	var cancelButton, receiveButton *widget.Button

	deleteAllFiles := func() {
		for fpath, fe := range fileentries {
			boxholder.Remove(fe)
			os.Remove(fpath)
			log.Tracef("Removed received file: %s", fpath)
			delete(fileentries, fpath)
		}
	}

	saveAllFiles := func() {
		if len(fileentries) == 0 {
			log.Error("no files to save")
			return
		}

		ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				log.Errorf("Error selecting folder: %v", err)
				return
			}
			if uri == nil {
				return
			}

			lastSaveDir = uri.Path()
			prog.Show()
			prog.Max = float64(len(fileentries))
			prog.SetValue(0)

			go func() {
				for fpath := range fileentries {
					dest := filepath.Join(lastSaveDir, filepath.Base(fpath))
					err := copyFile(fpath, dest)
					if err != nil {
						log.Errorf("Error saving file %s: %v", filepath.Base(fpath), err)
						continue
					}
					fyne.Do(func() {
						prog.SetValue(prog.Value + 1)
					})
				}
				fyne.Do(func() {
					prog.Hide()
					status.SetText(fmt.Sprintf("Saved all files to: %s", lastSaveDir))
				})
			}()
		}, w)
	}

<<<<<<< HEAD
				prog.Show()
				donechan := make(chan bool)
				var filename string
				receivednames := make(map[string]int)
				go func() {
					ticker := time.NewTicker(time.Millisecond * 100)
					for {
						select {
						case <-ticker.C:
							if receiver.Step2FileInfoTransferred {
								cnum := receiver.FilesToTransferCurrentNum
								fi := receiver.FilesToTransfer[cnum]
								filename = filepath.Base(fi.Name)
								receivednames[filename] = cnum
								topline.SetText(fmt.Sprintf("%s: %s (%d/%d)", lp("Receiving file"), filename, cnum+1, len(receiver.FilesToTransfer)))
								prog.Max = float64(fi.Size)
								prog.SetValue(float64(receiver.TotalSent))
							}
						case <-donechan:
							ticker.Stop()
							return
						}
					}
				}()
				cderr := os.Chdir(recvDir)
				if cderr != nil {
					log.Error("Unable to change to dir:", recvDir, cderr)
				}
				status.SetText("")
				rerr := receiver.Receive()
				donechan <- true
				prog.Hide()
				prog.SetValue(0)
				topline.SetText(lp("Enter code to download"))
				if rerr != nil {
					log.Error("Receive failed: " + rerr.Error())
				} else {
					filesReceived := make([]string, len(receivednames))
					var i int
					for f := range receivednames {
						filesReceived[i] = f
						i++
					}
					sort.Slice(filesReceived, func(i, j int) bool {
						return receivednames[filesReceived[i]] < receivednames[filesReceived[j]]
					})

					status.SetText(fmt.Sprintf("%s: %s", lp("Received"), strings.Join(filesReceived, ",")))
					filepath.Walk(recvDir, func(path string, info fs.FileInfo, err error) error {
						if err != nil {
							return err
						}
						if !info.IsDir() {
							var diagwg sync.WaitGroup
							diagwg.Add(1)
							savedialog := dialog.NewFileSave(func(f fyne.URIWriteCloser, e error) {
								var ofile io.WriteCloser
								var oerr error
								ofile = f
								oerr = e
								if oerr != nil {
									log.Error(oerr.Error())
									return
								}
								ifile, ierr := os.Open(path)
								if ierr != nil {
									log.Error(ierr.Error())
									return
								}
								nw, ew := io.Copy(ofile, ifile)
								if ew != nil {
									log.Errorf("%s: write failure from %s to %s, wrote %d bytes", ew.Error(), path, f.URI().String(), nw)
								}
								ifile.Close()
								ofile.Close()
								log.Tracef("saved (%s) to user path %s", path, f.URI().String())
								diagwg.Done()
							}, w)
							savedialog.SetFileName(filepath.Base(path))
							savedialog.Show()
							diagwg.Wait()
						}
						return nil
					})
				}
				// Clear recv dir after finished
				filepath.Walk(recvDir, func(path string, info fs.FileInfo, err error) error {
					if !info.IsDir() {
						os.Remove(path)
						log.Tracef("remove internal cache file %s", path)
					}
					return nil
				})
			}),
			prog,
			status,
			debugBox,
		))
=======
	resetReceiver := func() {
		prog.Hide()
		prog.SetValue(0)
		for _, obj := range activeButtonHolder.Objects {
			activeButtonHolder.Remove(obj)
		}
		activeButtonHolder.Add(receiveButton)

		topline.SetText(lp(""))
		recvEntry.Enable()
	}

	receiveButton = widget.NewButtonWithIcon(lp("Download"), theme.DownloadIcon(), func() {
		if recvEntry.Text == "" {
			log.Error("no receive code entered")
			dialog.ShowInformation(
				lp("Download"),
				lp("Enter code to download"),
				w,
			)
			return
		}
		if len(fileentries) > 0 {
			log.Error("save received files")
			dialog.ShowInformation(
				lp("Download"),
				lp("Save All"),
				w,
			)
			return
		}

		receiver, err := croc.New(croc.Options{
			IsSender:         false,
			SharedSecret:     recvEntry.Text,
			Debug:            crocDebugMode(),
			RelayAddress:     a.Preferences().String("relay-address"),
			RelayPorts:       strings.Split(a.Preferences().String("relay-ports"), ","),
			RelayPassword:    a.Preferences().String("relay-password"),
			Stdout:           false,
			NoPrompt:         true,
			DisableLocal:     a.Preferences().Bool("disable-local"),
			NoMultiplexing:   a.Preferences().Bool("disable-multiplexing"),
			OnlyLocal:        a.Preferences().Bool("force-local"),
			NoCompress:       a.Preferences().Bool("disable-compression"),
			Curve:            a.Preferences().String("pake-curve"),
			HashAlgorithm:    a.Preferences().String("croc-hash"),
			Overwrite:        true,
			ZipFolder:        false,
			GitIgnore:        false,
			MulticastAddress: a.Preferences().String("multicast-address"),
		})
		if err != nil {
			log.Errorf("Receive setup error: %s\n", err.Error())
			return
		}
		log.SetLevel(crocDebugLevel())
		log.Trace("croc receiver created")
		cderr := os.Chdir(recvDir)
		if cderr != nil {
			log.Error("Unable to change to dir:", recvDir, cderr)
		}
		log.Trace("cd", recvDir)

		var filename string
		status.SetText(fmt.Sprintf("%s: %s", lp("Receive Code"), recvEntry.Text))
		prog.Show()

		for _, obj := range activeButtonHolder.Objects {
			activeButtonHolder.Remove(obj)
		}
		activeButtonHolder.Add(cancelButton)

		donechan := make(chan bool)
		go func() {
			ticker := time.NewTicker(time.Millisecond * 100)
			for {
				select {
				case <-ticker.C:
					if receiver.Step2FileInfoTransferred {
						cnum := receiver.FilesToTransferCurrentNum
						fi := receiver.FilesToTransfer[cnum]
						filename = filepath.Base(fi.Name)
						fyne.Do(func() {
							topline.SetText(fmt.Sprintf("%s: %s(%d/%d)", lp("Receiving file"), filename, cnum+1, len(receiver.FilesToTransfer)))
							prog.Max = float64(fi.Size)
							prog.SetValue(float64(receiver.TotalSent))
						})
					}
				case <-donechan:
					ticker.Stop()
					return
				}
			}
		}()
>>>>>>> 7a60eeb (Finally, the Receive tab doesn't freeze anymore,)

		go func() {
			fyne.Do(recvEntry.Disable)
			ferr := receiver.Receive()
			donechan <- true
			if ferr != nil {
				log.Errorf("Receive failed: %s\n", ferr)
			} else {
				fyne.Do(func() {
					status.SetText(fmt.Sprintf("%s: %s", lp("Received"), filename))

					for _, fi := range receiver.FilesToTransfer {
						fpath := filepath.Join(recvDir, filepath.Base(fi.Name))
						labelFile := widget.NewLabel(filepath.Base(fpath))

						openButton := widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {
							ShowFileLocation(fpath, w)
						})

						deleteButton := widget.NewButtonWithIcon("", theme.CancelIcon(), func() {
							if fe, ok := fileentries[fpath]; ok {
								boxholder.Remove(fe)
								os.Remove(fpath)
								log.Tracef("Removed received file: %s", fpath)
								delete(fileentries, fpath)
							}
						})

						newentry := container.NewHBox(
							labelFile,
							layout.NewSpacer(),
							openButton,
							deleteButton,
						)
						fileentries[fpath] = newentry
						boxholder.Add(newentry)
					}
				})
			}
			fyne.Do(resetReceiver)
		}()

		go func() {
			select {
			case <-cancelchan:
				receiver.SuccessfulTransfer = true
				donechan <- true
				fyne.Do(func() {
					status.SetText(lp("Receive cancelled."))
				})
			}
			fyne.Do(resetReceiver)
		}()
	})

	cancelButton = widget.NewButtonWithIcon(lp("Cancel"), theme.CancelIcon(), func() {
		cancelchan <- true
	})

	activeButtonHolder.Add(receiveButton)

	deleteAllButton := widget.NewButtonWithIcon(lp("Delete All"), theme.DeleteIcon(), func() {
		dialog.ShowConfirm(lp("Delete all files"), lp("Are you sure you want to delete all received files?"), func(b bool) {
			if b {
				deleteAllFiles()
			}
		}, w)
	})

	saveAllButton := widget.NewButtonWithIcon(lp("Save All"), theme.FolderOpenIcon(), func() {
		saveAllFiles()
	})

	receiveTop := container.NewVBox(
		container.NewHBox(topline, layout.NewSpacer()),
		widget.NewForm(&widget.FormItem{Text: lp("Receive Code"), Widget: recvEntry}),
	)
	receiveBot := container.NewVBox(
		activeButtonHolder,
		prog,
		container.NewHBox(status),
		container.NewHBox(
			layout.NewSpacer(),
			saveAllButton,
			deleteAllButton,
		),
		debugBox,
	)

	return container.NewTabItemWithIcon(lp("Receive"), theme.DownloadIcon(),
		container.NewBorder(receiveTop, receiveBot, nil, nil, receiverScroller))
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func ShowFileLocation(path string, parent fyne.Window) {
	savedialog := dialog.NewFileSave(func(f fyne.URIWriteCloser, e error) {
		if f != nil {
			src, err := os.Open(path)
			if err != nil {
				log.Error(err)
				return
			}
			defer src.Close()

			_, err = io.Copy(f, src)
			if err != nil {
				log.Error(err)
			}
			f.Close()
		}
	}, parent)
	savedialog.SetFileName(filepath.Base(path))
	savedialog.Resize(parent.Canvas().Size())
	savedialog.Show()
}

// Big File Dialog
func ShowFolderOpen(callback func(fyne.ListableURI, error), parent fyne.Window) {
	switch runtime.GOOS {
	case "ios", "android":
		dialog.NewFolderOpen(callback, parent)
		return
	}
	fd := dialog.NewFolderOpen(callback, parent)
	fd.Resize(parent.Canvas().Size())
	fd.Show()
}
