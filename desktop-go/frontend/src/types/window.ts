declare global {
    interface Console {
        dev:(...args: any[]) => void;
    }

    interface Window {
        copyToClipboard: (text: string) => Promise<void>;
        copy: (value: any) => any;
    }
}