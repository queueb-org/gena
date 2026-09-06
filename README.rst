Gen-A
=====

Gen-A (Generators Assistant) is a small command-line helper for discovering
annotated Go packages and running Kubernetes code generators against them.

.. warning::

   Gen-A is alpha software. The command-line interface and generated-code
   workflow may change between ``0.0.x`` releases.

The initial ``v0.0.1`` release targets Linux on ``amd64`` and ``arm64``.
Windows, macOS, and other architectures are not guaranteed yet.

Installation
------------

Gen-A is distributed as a Go module during the alpha release series. Install
the selected version with Go 1.27 or newer:

.. code-block:: bash

   go install queueb.org/gena@v0.0.1
   gena version

Prebuilt binary archives are not published for ``v0.0.1``.

External generators
~~~~~~~~~~~~~~~~~~~

Gen-A does not install Kubernetes code generators. Install the binaries needed
by your workflow and make sure they are available through ``PATH`` before
running ``gena``.

The base ``--all`` set can be installed with:

.. code-block:: bash

   go install \
       k8s.io/code-generator/cmd/deepcopy-gen@latest \
       k8s.io/code-generator/cmd/defaulter-gen@latest \
       k8s.io/code-generator/cmd/conversion-gen@latest \
       k8s.io/code-generator/cmd/validation-gen@latest \
       k8s.io/code-generator/cmd/register-gen@latest

Client generation additionally requires:

.. code-block:: bash

   go install \
       k8s.io/code-generator/cmd/applyconfiguration-gen@latest \
       k8s.io/code-generator/cmd/client-gen@latest \
       k8s.io/code-generator/cmd/lister-gen@latest \
       k8s.io/code-generator/cmd/informer-gen@latest

Use a ``code-generator`` version compatible with the Kubernetes dependencies
of the target project.


Quick start
-----------

Add the required Kubernetes annotations to a package. For example, a
``doc.go`` file for deepcopy generation may start with:

.. code-block:: go

   // +k8s:deepcopy-gen=package
   package v1alpha1

Discover matching packages below the current module:

.. code-block:: bash

   gena tools discover .

The shorter ``tools`` alias is ``t``:

.. code-block:: bash

   gena t discover .

Run the five base generators:

.. code-block:: bash

   gena run --all --discover .

Client generators have contextual dependencies and should be run together in
this order:

.. code-block:: bash

   gena run \
       applyconfiguration-gen \
       client-gen \
       lister-gen \
       informer-gen \
       --discover .

The `example application <examples/app-1>`_ contains annotations and generated
output for both workflows.

Package discovery
-----------------

``gena tools discover`` recursively inspects one or more directory trees and
groups matching packages by their Go module. Without positional arguments it
uses the current directory.

Search several roots:

.. code-block:: bash

   #: or gena tools discover ./pkg ./internal ~/go/src/github.com/user/app
   gena t discover ./pkg ./internal ~/go/src/github.com/user/app

Search for another annotation:

.. code-block:: bash

   gena t discover \
       --annotation '+example:generate=true' \
       ./pkg

Write plain package names to a file:

.. code-block:: bash

   gena t discover \
       --output packages.txt \
       ./pkg ./internal

Render the grouped discovery result as JSON:

.. code-block:: bash

   gena t discover \
       --format='{{ json . }}' \
       ./pkg

The discovery command is primarily a diagnostic helper. Result ordering is not
part of its current interface.

Running generators
------------------

Run selected generators by name:

.. code-block:: bash

   gena run deepcopy-gen register-gen --discover .

``--all`` currently runs:

* ``deepcopy-gen``
* ``defaulter-gen``
* ``conversion-gen``
* ``validation-gen``
* ``register-gen``

Client generators are supported but are not included in ``--all``. They must
be requested explicitly in the order shown in the quick start.

By default, Gen-A changes into each discovered module before starting a
generator. With ``--change-dir=false``, generators run in the caller's current
working directory, so the caller must already be inside the target project.

Run ``gena run --help`` and ``gena tools discover --help`` for the complete
set of flags and corresponding environment variables.

Version information
-------------------

The application version comes from Go build information and the Git tag:

.. code-block:: bash

   gena version
   gena version --short

For the first release, the short command prints ``0.0.1``. A build made from
sources without VCS metadata reports ``devel`` and is not a release artifact.

Current limitations
-------------------

* Only Linux ``amd64`` and ``arm64`` are release targets.
* Generator binaries must be installed separately and available in ``PATH``.
* Prebuilt Gen-A archives are not provided during the alpha release series.
* ``openapi-gen``, ``pb-gen``, and other generators outside the documented base
  and client sets are not part of the ``v0.0.1`` support scope.
* Configuration and additional generation workflows remain under development.

Development
-----------

Install development tools and run the complete local verification:

.. code-block:: bash

   ./hack/install-dev-tools.sh
   ./hack/verify-local.sh

The project maintains 100% coverage for included code. Generated files and
other explicitly designated paths are intentionally excluded by
``hack/test.sh``.

Build a local binary with:

.. code-block:: bash

   ./hack/build.sh
   ./bin/gena version

License
-------

Gen-A is available under the `MIT License <LICENSE>`_.
