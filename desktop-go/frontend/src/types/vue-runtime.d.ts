import '@vue/runtime-core'

declare module '@vue/runtime-core' {
    interface ComponentCustomProperties {
        $toDate: (value: any) => Date;
        $isImage: (value: any) => boolean;
        $isVideo: (value: any) => boolean;
        $isAudio: (value: any) => boolean;
    }
}