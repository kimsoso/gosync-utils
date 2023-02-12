package command

// func (c *CommandBlk) sendCommands(wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	packNum := 0
// 	for {
// 		select {
// 		case <-c.ctx.Done():
// 			return
// 		case cmds := <-c.datachan:
// 			if IsStart(cmds) {
// 				if err := sonet.WriteCommandWithValue(c.conn, CBLOCKS_START, nil); err != nil {
// 					log.Println("send command blocks start error:", err)
// 					return
// 				} else {
// 					if _, cmd, err := sonet.Read(c.conn); err != nil {
// 						log.Println("send command blocks start feedback error:", err)
// 						return
// 					} else if cmd[0] != OK {
// 						log.Println("send command blocks start recieve wrong reply:", cmd[0])
// 						return
// 					}
// 				}

// 			} else if IsEnd(cmds) {
// 				if err := sonet.Write(c.conn, cmds); err != nil {
// 					log.Println("send command blocks fininsh error:", err)
// 				}
// 				return

// 			} else if err := sonet.SendPack(c.conn, uint32(packNum), cmds); err != nil {
// 				log.Println("send command blocks error:", err)
// 				return

// 			} else {
// 				if _packnum, cmd, err := sonet.Read(c.conn); err != nil {
// 					log.Println("rev command blocks error:", err)
// 					return
// 				} else if _packnum != uint32(packNum) {
// 					log.Println("send command blocks mismatch packnum:", packNum, _packnum)
// 					return
// 				} else if cmd[0] == CBLOCKS_TRANSFILE {
// 					offset := binary.BigEndian.Uint64(cmd[1:])
// 					path := c.conf.RealFilepath(utils.B2S(cmd[5:]))

// 					if _, err := sonet.SendFileDataRR(c.conn, path, int64(offset), 0, nil); err != nil {
// 						log.Println("send command blocks send file failed:", err)
// 						return
// 					}
// 				} else if cmd[0] != OK {
// 					log.Println("send command blocks recieve wrong reply:", cmd[0])
// 					return
// 				}

// 				packNum++
// 			}
// 		}
// 	}
// }
