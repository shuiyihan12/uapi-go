// Package usecase provides the use-case orchestration layer (Facade) for each domain.
//
// # Collapsible Layer
//
// Most Facades in this layer currently do only three things: fetch the service from the
// ServiceManager, ensure the trace ID, and forward the call to the services package. They
// carry almost no business logic of their own — these are "pass-through Facades".
//
// This layer is kept (rather than letting api call services directly) to balance
// "architectural consistency" against "simplicity" and to leave room for evolution:
//
//   - When a domain later needs input validation, enum normalization, cross-service
//     orchestration, or retry/fallback, the business logic naturally lands in this layer,
//     without changing the boundary between api and services;
//   - For purely pass-through domains, the current cost of this layer is just one function
//     call, which is negligible.
//
// Hence these Facades are marked as a "collapsible layer": in domains without business
// logic they are thin wrappers; as soon as a domain needs business semantics, this layer
// immediately earns its place. Follow the same pattern when adding new domains
// (NewXxxFacade + direct method forwarding), keep the layering uniform, and do not skip
// this layer for individual domains.
package usecase
