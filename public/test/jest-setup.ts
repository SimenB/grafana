// This import has side effects, and must be at the top so jQuery is made global first
import './global-jquery-shim';

import { TransformStream } from 'node:stream/web';
import { TextEncoder, TextDecoder } from 'util';

// we need to isolate the `@grafana/data` module here now that it depends on `@grafana/i18n`
jest.isolateModulesAsync(async () => {
  const { EventBusSrv } = await import('@grafana/data');
  const testAppEvents = new EventBusSrv();
  jest.mock('../app/core/core', () => ({
    ...jest.requireActual('../app/core/core'),
    appEvents: testAppEvents,
  }));
});
import { GrafanaBootConfig } from '@grafana/runtime';

import 'blob-polyfill';
import 'mutationobserver-shim';
import './mocks/workers';

import '../vendor/flot/jquery.flot';
import '../vendor/flot/jquery.flot.time';

const global = window as any;

// mock the default window.grafanaBootData settings
const settings: Partial<GrafanaBootConfig> = {
  featureToggles: {},
};
global.grafanaBootData = {
  settings,
  user: {
    locale: 'en-US',
  },
  navTree: [],
};

window.matchMedia = (query) => ({
  matches: false,
  media: query,
  onchange: null,
  addListener: jest.fn(), // Deprecated
  removeListener: jest.fn(), // Deprecated
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
  dispatchEvent: jest.fn(),
});

// mock the intersection observer and just say everything is in view
const mockIntersectionObserver = jest
  .fn()
  .mockImplementation((callback: (arg: IntersectionObserverEntry[]) => void) => ({
    observe: jest.fn().mockImplementation((elem: HTMLElement) => {
      callback([{ target: elem, isIntersecting: true }] as unknown as IntersectionObserverEntry[]);
    }),
    unobserve: jest.fn(),
    disconnect: jest.fn(),
  }));
global.IntersectionObserver = mockIntersectionObserver;
Object.defineProperty(document, 'fonts', {
  value: { ready: Promise.resolve({}) },
});

global.TextEncoder = TextEncoder;
global.TextDecoder = TextDecoder;
global.TransformStream = TransformStream;
// add scrollTo interface since it's not implemented in jsdom
Element.prototype.scrollTo = () => {};

const throwUnhandledRejections = () => {
  process.on('unhandledRejection', (err) => {
    throw err;
  });
};

throwUnhandledRejections();

// Used by useMeasure
global.ResizeObserver = class ResizeObserver {
  static #observationEntry: ResizeObserverEntry = {
    contentRect: {
      x: 1,
      y: 2,
      width: 500,
      height: 500,
      top: 100,
      bottom: 0,
      left: 100,
      right: 0,
    },
    target: {
      // Needed for react-virtual to work in tests
      getAttribute: () => 1,
    },
  } as unknown as ResizeObserverEntry;

  #isObserving = false;
  #callback: ResizeObserverCallback;

  constructor(callback: ResizeObserverCallback) {
    this.#callback = callback;
  }

  #emitObservation() {
    setTimeout(() => {
      if (!this.#isObserving) {
        return;
      }

      this.#callback([ResizeObserver.#observationEntry], this);
    });
  }

  observe() {
    this.#isObserving = true;
    this.#emitObservation();
  }

  disconnect() {
    this.#isObserving = false;
  }

  unobserve() {
    this.#isObserving = false;
  }
};

global.BroadcastChannel = class BroadcastChannel {
  onmessage() {}
  onmessageerror() {}
  postMessage(data: unknown) {}
  close() {}
  addEventListener() {}
  removeEventListener() {}
};
