import { connectStreamSource, disconnectStreamSource } from "@hotwired/turbo";

/**
 * Creates a persistent connection to an event source (SSE) for use with Turbo Streams
 */
export default class TurboStreamEventSource extends HTMLElement {
  get src() {
    return this.getAttribute("src");
  }

  set src(value) {
    if (value) {
      this.setAttribute("src", value);
    } else {
      this.removeAttribute("src");
    }
  }

  /**
   * Called when the element is inserted into the DOM.
   * Connects to Turbo Streams as a source and sets up the event source (SSE) connection
   * for streaming updates to turbo streams
   */
  async connectedCallback() {
    this.es = this.setupEventSource();
    connectStreamSource(this.es);
  }

  /**
   * Called when the element is removed from the DOM.
   * Disconnects from Turbo Streams and deletes the WebSocket
   */
  disconnectedCallback() {
    this.es.close();
    disconnectStreamSource(this.es);
    if (this.es) {
      this.es = null;
    }
  }

  /**
   * Called in response to a websocket message. Unpacks the websocket message
   * and dispatches it as a new MessageEvent to Turbo Streams.
   * 
   * @param {MessageEvent} messageEvent The original message to dispatch
   */
  dispatchMessageEvent(messageEvent) {
    const event = new MessageEvent("message", { data: messageEvent.data });
    this.dispatchEvent(event);
  }

  /**
   * Creates a WebSocket by combining the window host and the element's src attribute and wires it to dispatch messages
   */
  setupEventSource() {
    let es = new EventSource(`http://${window.location.host}${this.src}`);
    es.onmessage = (msg) => {
      this.dispatchMessageEvent(msg);
    };
    return es;
  }
}
