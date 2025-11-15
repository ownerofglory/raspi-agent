import {useContext, useState} from "react";
import AuthCtx from "../../context/auth.ts";

const backendUrl = "http://localhost:8000/raspi-agent/api/v1/auth/login";

/**
 * Expected response format from the backend login API.
 */
interface LoginResult {
    id: string;
    token: string;
}

/**
 * Login page component.
 *
 * Handles user sign-in, sends credentials to the backend,
 * and stores the returned token/id inside the global AuthContext.
 */
export default function LoginPage() {
    /** User's email/username input */
    const [email, setEmail] = useState<string>();

    /** User's password input */
    const [password, setPassword] = useState<string>();

    /** Auth context used to store login results */
    const authCtx = useContext(AuthCtx);

    /**
     * Triggered when the user clicks the "Sign In" button.
     *
     * Sends a POST request to the backend with email/password.
     * On success, updates global auth state via AuthContext.
     */
    const onLogin = () => {
        const localLogin = {email, password}

        fetch(backendUrl, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(localLogin),
        }).then<LoginResult>(res => res.json())
        .then(res => authCtx?.setAuth(res))
    }

    return (
        <div className="font-display bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark min-h-screen flex items-center justify-center p-4">
            <div className="w-full max-w-md rounded-xl border border-border-light dark:border-border-dark bg-card-light dark:bg-card-dark shadow-lg">
                <div className="flex flex-col p-8">
                    {/* Header */}
                    <div className="mb-8 flex flex-col items-center gap-4 text-center">
                        <div className="h-12 w-12 text-primary">
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
                        <h1 className="text-2xl font-bold">DIY Voice Assistant</h1>
                        <p className="text-text-muted-light dark:text-text-muted-dark">
                            Welcome back! Please sign in to continue.
                        </p>
                    </div>

                    {/* Login Form */}
                    <form action="#" method="POST" className="flex flex-col gap-4">
                        <div>
                            <label
                                htmlFor="username"
                                className="mb-1 block text-sm font-medium"
                            >
                                Username
                            </label>
                            <input
                                id="username"
                                name="username"
                                type="text"
                                placeholder="Enter your username"
                                required
                                className="w-full rounded-lg border-border-light bg-background-light p-3 placeholder-text-muted-light
                focus:border-primary focus:ring-2 focus:ring-primary/30
                dark:border-border-dark dark:bg-background-dark dark:placeholder-text-muted-dark dark:focus:border-primary"
                                onInput={e => setEmail((e.target as HTMLInputElement).value)}
                            />
                        </div>

                        <div>
                            <div className="mb-1 flex items-center justify-between">
                                <label
                                    htmlFor="password"
                                    className="text-sm font-medium"
                                >
                                    Password
                                </label>
                                <a
                                    href="#"
                                    className="text-sm font-medium text-primary hover:underline"
                                >
                                    Forgot password?
                                </a>
                            </div>
                            <input
                                id="password"
                                name="password"
                                type="password"
                                placeholder="Enter your password"
                                required
                                className="w-full rounded-lg border-border-light bg-background-light p-3 placeholder-text-muted-light
                focus:border-primary focus:ring-2 focus:ring-primary/30
                dark:border-border-dark dark:bg-background-dark dark:placeholder-text-muted-dark dark:focus:border-primary"
                                onInput={e => setPassword((e.target as HTMLInputElement).value)}
                            />
                        </div>

                        <button
                            type="submit"
                            className="mt-4 flex h-12 w-full cursor-pointer items-center justify-center gap-2 overflow-hidden
              rounded-lg bg-primary px-6 text-base font-bold text-white shadow-lg shadow-primary/30
              transition-all hover:bg-primary/90"
                            onClick={() => onLogin()}
                        >
                            <span className="truncate">Sign In</span>
                        </button>
                    </form>

                    {/* Divider */}
                    <div className="my-6 flex items-center">
                        <div className="h-px flex-grow bg-border-light dark:bg-border-dark"></div>
                        <span className="mx-4 text-sm text-text-muted-light dark:text-text-muted-dark">
              OR
            </span>
                        <div className="h-px flex-grow bg-border-light dark:bg-border-dark"></div>
                    </div>

                    {/* Google Button */}
                    <button
                        type="button"
                        className="flex h-12 w-full cursor-pointer items-center justify-center gap-3 overflow-hidden
            rounded-lg border border-border-light bg-card-light px-6 text-base font-bold text-text-light shadow-sm
            transition-all hover:bg-black/5 dark:border-border-dark dark:bg-card-dark dark:text-text-dark dark:hover:bg-white/5"
                    >

                        <span className="truncate">Sign in with Google</span>
                    </button>
                </div>
            </div>
        </div>
    );
}
