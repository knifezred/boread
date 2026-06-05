declare module '@ugreen-nas/core' {
  const UGOSCore: {
    init(): void;
  };
  export default UGOSCore;
}

declare module '@ugreen-nas/core/cloudWindow' {
  interface CloudWindow {
    winId: string;
    useCapacity(capName: string, data?: unknown, timeout?: number): Promise<Record<string, unknown>>;
    on(event: string, callback: (...args: unknown[]) => void): void;
  }
  const cloudWindow: CloudWindow;
  export default cloudWindow;
}
