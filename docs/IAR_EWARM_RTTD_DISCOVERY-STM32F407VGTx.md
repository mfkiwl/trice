# IAR_EWARM_RTTD_DISCOVERY-STM32F407VGTx
- direct RTT
- About: [SeggerRTT.md](SeggerRTT.md)

- Use "C:\ST\STLinkReflash_190812_Segger\STLinkReflash.exe" to convert On-board ST-Link to On-board J-Link (optional and you can undo this step)
- Start ST CubeMX
  - select DISCOVERY_STM32F407VGTx
  - generate for EWARM
- Open EWARM project & compile
  - Add Group SeggerRTT
  - Add all files from C:\repos\trice\third_party\SEGGER_RTT_V672b\RTT - also SEGGER_RTT_ASM_ARMv7M.S
  - Add group Trice
  - Add files from C:\repos\trice\srcTrice.C but not xteaCrypto.* and not triceWrite.c
  - Add TRICE_VARIANT=SEGGER_RTTD in project settings
  - Add include paths for RTT & Trice
  - Add Build action "trice u -src $PROJ_DIR$/.. -src $PROJ_DIR$/../../srcTrice.C -list $PROJ_DIR$/../../til.json -v"
  - Determine J-Link identifier (for example with J-Link-Mem)
  - Select J-Link as debugger
  - Switch J-Link to SWD
- ~ 180 clocks

![IAR_EWARM_RTTD_DISCOVERY-STM32F407VGTx_0.PNG](./README.media/IAR_EWARM_RTTD_DISCOVERY-STM32F407VGTx_0.PNG)
- ~ 80 clocks
  
![IAR_EWARM_RTTD_DISCOVERY-STM32F407VGTx_1.PNG](./README.media/IAR_EWARM_RTTD_DISCOVERY-STM32F407VGTx_1.PNG)

- ~100 clocks

![IAR_EWARM_RTTD_DISCOVERY-STM32F407VGTx.PNG](./README.media/IAR_EWARM_RTTD_DISCOVERY-STM32F407VGTx.PNG)