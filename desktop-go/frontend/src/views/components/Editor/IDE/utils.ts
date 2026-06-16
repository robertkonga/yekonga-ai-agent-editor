



// UI color badge selector based on language tag
export const getFileColorClass = (lang: string) => {
    switch (lang) {
        case 'javascript': return 'text-yellow-400';
        case 'css': return 'text-blue-400';
        case 'json': return 'text-emerald-400';
        default: return getIconColor(lang);
    }
};

export const getIconColor = (filename: string): string => {
    const ext = filename.substring(filename.lastIndexOf('.')).toLowerCase();
    switch (ext) {
        // Systems & Backend
        case '.go': return 'text-cyan-400';
        case '.rs': return 'text-orange-600';
        case '.py': return 'text-sky-500';
        case '.java': return 'text-red-500';
        case '.cpp':
        case '.hpp':
        case '.h': return 'text-blue-600';
        case '.cs': return 'text-purple-600';

        // Scripts & Infrastructure
        case '.sh':
        case '.bash': return 'text-emerald-500';
        case '.yaml':
        case '.yml':
        case '.toml': return 'text-indigo-500';
        case '.sql': return 'text-pink-500';

        // Core Web
        case '.ts':
        case '.tsx': return 'text-blue-500';
        case '.js':
        case '.jsx': return 'text-yellow-400';
        case '.vue': return 'text-emerald-400';
        case '.html': return 'text-orange-500';
        case '.css':
        case '.scss': return 'text-teal-400';
        case '.json': return 'text-amber-400';
        case '.md': return 'text-slate-400';

        default: return 'text-slate-600';
    }
};