import type { Workspace } from "./store/workspace";

const WORKSPACE = "__WORKSPACES_DATA__"

window.copyToClipboard = async (text: string): Promise<void> => {
    if (navigator.clipboard && window.isSecureContext) {
        return navigator.clipboard.writeText(text);
    } else {
        let textArea = document.createElement("textarea");
        textArea.value = text;
        textArea.style.position = "fixed";  // Avoid scrolling to bottom
        textArea.style.left = "-999999px";
        textArea.style.top = "-999999px";
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();

        return new Promise((res, rej) => {
            document.execCommand('copy') ? res() : rej();
            textArea.remove();
        });
    }
}

window.copy = (value: any) => {
    return JSON.parse(JSON.stringify(value));
}

window.toDate = (value: any): Date => {
    let date: Date = new Date();

    if(typeof value == 'string') {
        try {
            date = new Date(Date.parse(value));
        } catch (error) {
            console.log(error)
        }
    } else if(value instanceof Date) {
        date = value;
    }

    return date;
}

window.dateToNumber = (value: any): number => {
    let date = window.toDate(value);
    
    return date.getTime();
}

window.savaLocalData = (data: Record<string, Workspace>): void => {
    data = window.copy(data);

    for (const key in data) {
        if (!Object.hasOwn(data, key)) continue;

        data[key].workspaceFiles = [];
    }

    window.localStorage.setItem(WORKSPACE, JSON.stringify(data));
} 

window.fetchLocalData = (): Record<string, Workspace> => {
    let data: Record<string, Workspace> = {};

    try {
        let dataString = window.localStorage.getItem(WORKSPACE);
        let value = JSON.parse(dataString || "{}");

        if(value) {
            data = value;
        }
    } catch (error) {
        console.log(error);
    }

    return data
} 

window.isImage = (url: any): boolean => {
    if (typeof url !== 'string') return false;
    return /\.(jpg|jpeg|png|webp|avif|gif|svg|bmp|ico)(\?.*)?$/i.test(url);
};

window.isVideo = (url: any): boolean => {
    if (typeof url !== 'string') return false;
    return /\.(mp4|webm|ogg|avi|mov|wmv|flv|mkv)(\?.*)?$/i.test(url);
};

window.isAudio = (url: any): boolean => {
    if (typeof url !== 'string') return false;
    return /\.(mp3|wav|ogg|flac|aac|wma)(\?.*)?$/i.test(url);
};

