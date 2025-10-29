package drm

type drmIrqBusID struct {
	irq     int32 /**< IRQ number */
	busnum  int32 /**< bus number */
	devnum  int32 /**< device number */
	funcnum int32 /**< function number */
}
