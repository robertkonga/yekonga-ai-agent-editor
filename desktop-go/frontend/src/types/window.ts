import type { Workspace } from "@/store/workspace";

declare global {
    interface Console {
        dev:(...args: any[]) => void;
    }

    interface Window {
        copyToClipboard: (text: string) => Promise<void>;
        copy: (value: any) => any;
        savaLocalData:  (data: Record<string, Workspace>) => void;
        fetchLocalData: () => Record<string, Workspace>;
        toDate: (value: any) => Date;
        dateToNumber: (value: any) => number;
    }
}