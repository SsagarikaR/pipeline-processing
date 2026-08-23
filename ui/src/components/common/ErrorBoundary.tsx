import  { Component, ErrorInfo } from 'react';
import toast from 'react-hot-toast';

import { ErrorBoundaryProps, ErrorBoundaryState } from '../../types/common';
import { ERROR_BOUNDARY_TEXTS } from '../../constants/common';

/**
 * Catches render errors anywhere below it in the component tree and
 * shows a fallback "Something went wrong" screen with a reload button,
 * instead of letting the whole app crash to a blank page.
 */
class ErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  public state: ErrorBoundaryState = {
    hasError: false
  };

  /** React calls this after a child throws, to switch into the fallback UI. */
  public static getDerivedStateFromError(_: Error): ErrorBoundaryState {
    return { hasError: true };
  }

  /** Logs the error for debugging and shows the user a toast about it. */
  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Uncaught error:', error, errorInfo);
    toast.error(ERROR_BOUNDARY_TEXTS.toastMessage(error.message));
  }

  public render() {
    if (this.state.hasError) {
      return (
        <div className="flex flex-col items-center justify-center min-h-screen bg-neutral-50 text-neutral-800">
          <h1 className="text-2xl font-bold mb-4">{ERROR_BOUNDARY_TEXTS.TITLE}</h1>
          <p className="text-neutral-600 mb-6">{ERROR_BOUNDARY_TEXTS.MESSAGE}</p>
          <button
            className="px-4 py-2 bg-brand-600 text-white rounded hover:bg-brand-700 transition"
            onClick={() => window.location.reload()}
          >
            {ERROR_BOUNDARY_TEXTS.RELOAD_BUTTON}
          </button>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
