export default function LandingPage() {
    return (
        <div className="font-display bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark min-h-screen flex flex-col relative">
            {/* Header */}
            <header className="absolute top-0 z-10 flex w-full items-center justify-between px-4 py-4 sm:px-6 md:px-8 lg:px-10">
                <div className="flex items-center gap-3">
                    <div className="h-8 w-8 text-primary">
                        <svg
                            fill="none"
                            viewBox="0 0 48 48"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <g clipPath="url(#clip0_6_319)">
                                <path
                                    d="M8.57829 8.57829C5.52816 11.6284 3.451 15.5145 2.60947 19.7452C1.76794 23.9758 2.19984 28.361 3.85056 32.3462C5.50128 36.3314 8.29667 39.7376 11.8832 42.134C15.4698 44.5305 19.6865 45.8096 24 45.8096C28.3135 45.8096 32.5302 44.5305 36.1168 42.134C39.7033 39.7375 42.4987 36.3314 44.1494 32.3462C45.8002 28.361 46.2321 23.9758 45.3905 19.7452C44.549 15.5145 42.4718 11.6284 39.4217 8.57829L24 24L8.57829 8.57829Z"
                                    fill="currentColor"
                                ></path>
                            </g>
                            <defs>
                                <clipPath id="clip0_6_319">
                                    <rect fill="white" height="48" width="48"></rect>
                                </clipPath>
                            </defs>
                        </svg>
                    </div>
                    <h2 className="text-xl font-bold">DIY Voice Assistant</h2>
                </div>

                <div className="flex items-center gap-2 sm:gap-4">
                    <a
                        href="#"
                        className="flex h-10 cursor-pointer items-center justify-center rounded-lg px-4 text-sm font-bold text-text-light dark:text-text-dark hover:bg-black/5 dark:hover:bg-white/5 sm:text-base"
                    >
                        Login
                    </a>
                    <a
                        href="#"
                        className="flex h-10 cursor-pointer items-center justify-center rounded-lg bg-primary px-4 text-sm font-bold text-white shadow-lg shadow-primary/30 transition-all hover:bg-primary/90 sm:text-base"
                    >
                        Register
                    </a>
                </div>
            </header>

            {/* Main Content */}
            <main className="flex flex-1 items-center justify-center px-4 py-20 text-center">
                <div className="flex w-full max-w-3xl flex-col items-center gap-8">
                    <div className="flex flex-col gap-4">
                        <h1 className="text-4xl font-bold tracking-tighter sm:text-5xl md:text-6xl">
                            Build Your Own Private Voice Assistant
                        </h1>
                        <p className="mx-auto max-w-xl text-base text-text-muted-light dark:text-text-muted-dark sm:text-lg md:text-xl">
                            Transform your Raspberry Pi into a powerful, privacy-focused smart
                            assistant. Control your smart home, get answers, and more, all
                            running locally on your own device.
                        </p>
                    </div>

                    {/* Buttons */}
                    <div className="flex flex-col items-center gap-4 sm:flex-row">
                        <a
                            href="#"
                            className="flex h-12 w-full min-w-[180px] cursor-pointer items-center justify-center gap-2 overflow-hidden rounded-lg bg-primary px-6 text-base font-bold text-white shadow-lg shadow-primary/30 transition-all hover:bg-primary/90 sm:w-auto"
                        >
                            <span className="material-symbols-outlined">rocket_launch</span>
                            <span className="truncate">Get Started Now</span>
                        </a>

                        <a
                            href="#"
                            className="flex h-12 w-full min-w-[180px] cursor-pointer items-center justify-center gap-2 overflow-hidden rounded-lg bg-black/5 px-6 text-base font-bold text-text-light dark:bg-white/10 dark:text-text-dark hover:bg-black/10 dark:hover:bg-white/20 sm:w-auto"
                        >
                            <span className="material-symbols-outlined">library_books</span>
                            <span className="truncate">Read the Docs</span>
                        </a>
                    </div>

                    {/* Feature Cards */}
                    <div className="mt-8 grid grid-cols-1 gap-6 text-left sm:grid-cols-2 md:grid-cols-3">
                        <FeatureCard
                            icon="mic"
                            title="Custom Wake Word"
                            description="Train and use your own unique wake word for a truly personal experience."
                        />
                        <FeatureCard
                            icon="graphic_eq"
                            title="Speech-to-Text"
                            description="On-device or cloud-based transcription to understand your commands accurately."
                        />
                        <FeatureCard
                            icon="auto_awesome"
                            title="LLM Reasoning"
                            description="Leverage local or API-based Large Language Models for intelligent responses."
                        />
                    </div>
                </div>
            </main>
        </div>
    );
}

// Reusable feature card component
function FeatureCard({ icon, title, description }: {icon: string, title: string, description: string}) {
    return (
        <div className="flex flex-col gap-2 rounded-xl border border-border-light dark:border-border-dark bg-card-light dark:bg-card-dark p-4">
            <span className="material-symbols-outlined text-primary">{icon}</span>
            <h3 className="font-bold">{title}</h3>
            <p className="text-sm text-text-muted-light dark:text-text-muted-dark">
                {description}
            </p>
        </div>
    );
}

